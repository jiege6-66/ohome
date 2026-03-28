package service

import (
	"context"
	"errors"
	"ohome/global"
	"ohome/model"
	"sync"

	"github.com/spf13/viper"
)

const (
	defaultQuarkTransferTaskWorkers   = 2
	defaultQuarkTransferTaskQueueSize = 32
	quarkTransferTaskQueueFullMessage = "转存任务队列已满，请稍后重试"
)

var errQuarkTransferTaskQueueFull = errors.New(quarkTransferTaskQueueFullMessage)

type quarkTransferExecuteFunc func(ctx context.Context, task quarkSaveTask) (quarkSaveResult, error)

type quarkTransferTaskExecutorConfig struct {
	Workers   int
	QueueSize int
}

type quarkTransferTaskExecution struct {
	task     model.QuarkTransferTask
	saveTask quarkSaveTask
}

type quarkTransferTaskExecutor struct {
	jobs    chan quarkTransferTaskExecution
	execute quarkTransferExecuteFunc
}

var (
	quarkTransferTaskExecutorMu sync.Mutex
	activeQuarkTransferExecutor *quarkTransferTaskExecutor
	quarkTransferTaskExecute    quarkTransferExecuteFunc = executeQuarkTransfer
)

func EnsureQuarkTransferTaskExecutor() {
	getOrCreateQuarkTransferTaskExecutor()
}

func getOrCreateQuarkTransferTaskExecutor() *quarkTransferTaskExecutor {
	quarkTransferTaskExecutorMu.Lock()
	defer quarkTransferTaskExecutorMu.Unlock()

	if activeQuarkTransferExecutor != nil {
		return activeQuarkTransferExecutor
	}

	config := loadQuarkTransferTaskExecutorConfig()
	activeQuarkTransferExecutor = newQuarkTransferTaskExecutor(config, quarkTransferTaskExecute)
	if global.Logger != nil {
		global.Logger.Infof(
			"Quark Transfer Executor Started: workers=%d queueSize=%d",
			config.Workers,
			config.QueueSize,
		)
	}
	return activeQuarkTransferExecutor
}

func loadQuarkTransferTaskExecutorConfig() quarkTransferTaskExecutorConfig {
	workers := viper.GetInt("quark.transfer.workers")
	queueSize := viper.GetInt("quark.transfer.queueSize")

	if workers <= 0 {
		workers = defaultQuarkTransferTaskWorkers
	}
	if queueSize <= 0 {
		queueSize = defaultQuarkTransferTaskQueueSize
	}
	if queueSize < workers {
		queueSize = workers
	}

	return quarkTransferTaskExecutorConfig{
		Workers:   workers,
		QueueSize: queueSize,
	}
}

func newQuarkTransferTaskExecutor(
	config quarkTransferTaskExecutorConfig,
	execute quarkTransferExecuteFunc,
) *quarkTransferTaskExecutor {
	if execute == nil {
		execute = executeQuarkTransfer
	}

	executor := &quarkTransferTaskExecutor{
		jobs:    make(chan quarkTransferTaskExecution, config.QueueSize),
		execute: execute,
	}
	for i := 0; i < config.Workers; i++ {
		go executor.worker()
	}
	return executor
}

func (e *quarkTransferTaskExecutor) enqueue(job quarkTransferTaskExecution) error {
	select {
	case e.jobs <- job:
		return nil
	default:
		return errQuarkTransferTaskQueueFull
	}
}

func (e *quarkTransferTaskExecutor) worker() {
	service := &QuarkTransferTaskService{}
	for job := range e.jobs {
		service.runTransferTask(job.task, job.saveTask, e.execute)
	}
}

func (e *quarkTransferTaskExecutor) shutdown() {
	close(e.jobs)
}

func replaceQuarkTransferTaskExecutorForTest(
	config quarkTransferTaskExecutorConfig,
	execute quarkTransferExecuteFunc,
) func() {
	executor := newQuarkTransferTaskExecutor(config, execute)

	quarkTransferTaskExecutorMu.Lock()
	previous := activeQuarkTransferExecutor
	activeQuarkTransferExecutor = executor
	quarkTransferTaskExecutorMu.Unlock()

	return func() {
		quarkTransferTaskExecutorMu.Lock()
		activeQuarkTransferExecutor = previous
		quarkTransferTaskExecutorMu.Unlock()
		executor.shutdown()
	}
}

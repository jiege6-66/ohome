package service

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"ohome/global"
	"ohome/model"
)

func TestQuarkTransferTaskExecutorRespectsConcurrencyLimit(t *testing.T) {
	restoreDB := setupPhase2TestDB(t)
	defer restoreDB()

	if err := global.DB.AutoMigrate(&model.AppMessage{}); err != nil {
		t.Fatalf("AutoMigrate(AppMessage) error = %v", err)
	}

	release := make(chan struct{})
	started := make(chan uint, 3)
	var running int32
	var maxRunning int32

	restoreExecutor := replaceQuarkTransferTaskExecutorForTest(
		quarkTransferTaskExecutorConfig{Workers: 2, QueueSize: 3},
		func(ctx context.Context, task quarkSaveTask) (quarkSaveResult, error) {
			current := atomic.AddInt32(&running, 1)
			for {
				previous := atomic.LoadInt32(&maxRunning)
				if current <= previous || atomic.CompareAndSwapInt32(&maxRunning, previous, current) {
					break
				}
			}
			defer atomic.AddInt32(&running, -1)

			started <- task.ID
			select {
			case <-release:
				return quarkSaveResult{Status: "ok", Message: "done", SavedCount: 1}, nil
			case <-ctx.Done():
				return quarkSaveResult{Status: "fail", Message: ctx.Err().Error()}, ctx.Err()
			}
		},
	)
	defer restoreExecutor()

	service := &QuarkTransferTaskService{}
	submitted := make([]model.QuarkTransferTask, 0, 3)
	for i := 0; i < 3; i++ {
		task, err := service.submit(quarkTransferTaskSubmission{
			displayName: "task",
			shareURL:    "https://example.com/share",
			savePath:    "/library",
			sourceType:  model.QuarkTransferTaskSourceSearchManual,
			ownerUserID: 1,
		})
		if err != nil {
			t.Fatalf("submit(task=%d) error = %v", i+1, err)
		}
		submitted = append(submitted, task)
	}

	for i := 0; i < 2; i++ {
		select {
		case <-started:
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for worker start")
		}
	}

	select {
	case taskID := <-started:
		t.Fatalf("task #%d started before a worker was released", taskID)
	case <-time.After(200 * time.Millisecond):
	}

	queuedTask := loadTransferTaskByID(t, submitted[2].ID)
	if queuedTask.Status != model.QuarkTransferTaskStatusQueued {
		t.Fatalf("queued task status = %s, want %s", queuedTask.Status, model.QuarkTransferTaskStatusQueued)
	}
	if queuedTask.StartedAt != nil {
		t.Fatalf("queued task startedAt = %v, want nil", queuedTask.StartedAt)
	}

	if max := atomic.LoadInt32(&maxRunning); max != 2 {
		t.Fatalf("max running = %d, want 2", max)
	}

	close(release)

	for _, item := range submitted {
		waitForTransferTaskStatus(t, item.ID, model.QuarkTransferTaskStatusSuccess)
	}
}

func TestQuarkTransferTaskQueueFullFailsFast(t *testing.T) {
	restoreDB := setupPhase2TestDB(t)
	defer restoreDB()

	if err := global.DB.AutoMigrate(&model.AppMessage{}); err != nil {
		t.Fatalf("AutoMigrate(AppMessage) error = %v", err)
	}

	release := make(chan struct{})
	started := make(chan struct{}, 1)

	restoreExecutor := replaceQuarkTransferTaskExecutorForTest(
		quarkTransferTaskExecutorConfig{Workers: 1, QueueSize: 1},
		func(ctx context.Context, task quarkSaveTask) (quarkSaveResult, error) {
			select {
			case started <- struct{}{}:
			default:
			}
			select {
			case <-release:
				return quarkSaveResult{Status: "ok", Message: "done", SavedCount: 1}, nil
			case <-ctx.Done():
				return quarkSaveResult{Status: "fail", Message: ctx.Err().Error()}, ctx.Err()
			}
		},
	)
	defer restoreExecutor()

	service := &QuarkTransferTaskService{}
	first, err := service.submit(quarkTransferTaskSubmission{
		displayName: "first",
		shareURL:    "https://example.com/1",
		savePath:    "/library",
		sourceType:  model.QuarkTransferTaskSourceSearchManual,
		ownerUserID: 1,
	})
	if err != nil {
		t.Fatalf("submit(first) error = %v", err)
	}

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for first task to start")
	}

	second, err := service.submit(quarkTransferTaskSubmission{
		displayName: "second",
		shareURL:    "https://example.com/2",
		savePath:    "/library",
		sourceType:  model.QuarkTransferTaskSourceSearchManual,
		ownerUserID: 1,
	})
	if err != nil {
		t.Fatalf("submit(second) error = %v", err)
	}

	if _, err := service.submit(quarkTransferTaskSubmission{
		displayName: "overflow",
		shareURL:    "https://example.com/3",
		savePath:    "/library",
		sourceType:  model.QuarkTransferTaskSourceSearchManual,
		ownerUserID: 1,
	}); err == nil || err.Error() != quarkTransferTaskQueueFullMessage {
		t.Fatalf("submit(overflow) error = %v, want %s", err, quarkTransferTaskQueueFullMessage)
	}

	queuedTask := loadTransferTaskByID(t, second.ID)
	if queuedTask.Status != model.QuarkTransferTaskStatusQueued {
		t.Fatalf("second task status = %s, want %s", queuedTask.Status, model.QuarkTransferTaskStatusQueued)
	}

	var overflowTask model.QuarkTransferTask
	if err := global.DB.Where("display_name = ?", "overflow").First(&overflowTask).Error; err != nil {
		t.Fatalf("load overflow task error = %v", err)
	}
	if overflowTask.Status != model.QuarkTransferTaskStatusFailed {
		t.Fatalf("overflow task status = %s, want %s", overflowTask.Status, model.QuarkTransferTaskStatusFailed)
	}
	if overflowTask.ResultMessage != quarkTransferTaskQueueFullMessage {
		t.Fatalf("overflow task message = %s, want %s", overflowTask.ResultMessage, quarkTransferTaskQueueFullMessage)
	}

	close(release)
	waitForTransferTaskStatus(t, first.ID, model.QuarkTransferTaskStatusSuccess)
	waitForTransferTaskStatus(t, second.ID, model.QuarkTransferTaskStatusSuccess)
}

func loadTransferTaskByID(t *testing.T, id uint) model.QuarkTransferTask {
	t.Helper()

	var task model.QuarkTransferTask
	if err := global.DB.First(&task, id).Error; err != nil {
		t.Fatalf("load transfer task #%d error = %v", id, err)
	}
	return task
}

func waitForTransferTaskStatus(t *testing.T, id uint, want string) {
	t.Helper()

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		task := loadTransferTaskByID(t, id)
		if task.Status == want {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}

	task := loadTransferTaskByID(t, id)
	t.Fatalf("task #%d status = %s, want %s", id, task.Status, want)
}

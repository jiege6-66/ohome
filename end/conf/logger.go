package conf

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() *zap.SugaredLogger {
	// Warn level and above only
	level := zapcore.WarnLevel

	core := zapcore.NewCore(getEncoder(), zapcore.NewMultiWriteSyncer(getWriteSync(), zapcore.AddSync(os.Stdout)), level)
	return zap.New(core).Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Local().Format(time.DateTime))
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}
func getWriteSync() zapcore.WriteSyncer {
	//获取分割符
	separator := string(filepath.Separator)
	//获取项目根目录
	stRootDir, _ := os.Getwd()
	// 生成log的目录
	logPath := stRootDir + separator + "log" + separator + time.Now().Format(time.DateOnly) + ".log"

	lumberjackSyncer := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    viper.GetInt("logger.MaxSize"),    // 日志切割的开始
		MaxBackups: viper.GetInt("logger.MaxBackups"), //保留的最大数量
		MaxAge:     viper.GetInt("logger.MaxAge"),     //保留的最长时间
		Compress:   false,                             // disabled by default
	}

	return zapcore.AddSync(lumberjackSyncer)
}

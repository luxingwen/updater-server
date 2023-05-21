package logger

import (
	"updater-server/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*zap.SugaredLogger
}

func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{l.SugaredLogger.With(args...)}
}

func NewLogger(config config.LogConfig) *Logger {
	writeSyncer := getLogWriter(config)
	encoder := getEncoder(config.Format)

	var logLevel zapcore.Level
	err := logLevel.Set(config.Level)
	if err != nil {
		logLevel = zap.InfoLevel
	}

	core := zapcore.NewCore(encoder, writeSyncer, logLevel)

	logger := zap.New(core, zap.AddCaller())

	return &Logger{logger.Sugar()}
}

func getEncoder(format string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder

	switch format {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		return zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}

func getLogWriter(config config.LogConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize, // megabytes
		MaxBackups: 5,
		MaxAge:     config.MaxAge,   //days
		Compress:   config.Compress, // disabled by default
		LocalTime:  true,
	}

	return zapcore.AddSync(lumberJackLogger)
}

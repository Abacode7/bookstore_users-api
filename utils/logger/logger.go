package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

var (
	logger *zap.Logger
)

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{"stdout"},
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	var err error
	logger, err = logConfig.Build()
	if err != nil {
		log.Fatalln(err)
	}
}

/// Info prints an info level log with message
func Info(message string, tag ...zap.Field) {
	logger.Info(message, tag...)
	logger.Sync()
}

/// Error prints an error level log with message and error
func Error(message string, err error, tag ...zap.Field) {
	zapErr := zap.NamedError("error", err)
	tag = append(tag, zapErr)
	logger.Error(message, tag...)
	logger.Sync()
}

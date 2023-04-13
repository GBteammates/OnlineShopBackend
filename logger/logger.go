package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Logger *zap.Logger
}

func NewLogger(logLevel string) *Logger {
	atomicLevel := zap.NewAtomicLevel()

	switch logLevel {
	case "info":
		{
			atomicLevel.SetLevel(zap.InfoLevel)
		}
	case "warning":
		{
			atomicLevel.SetLevel(zap.WarnLevel)
		}
	case "debug":
		{
			atomicLevel.SetLevel(zap.DebugLevel)
		}
	case "error":
		{
			atomicLevel.SetLevel(zap.ErrorLevel)

		}
	case "panic":
		{
			atomicLevel.SetLevel(zap.PanicLevel)
		}
	case "fatal":
		{
			atomicLevel.SetLevel(zap.FatalLevel)
		}
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atomicLevel,
	), zap.AddCaller())

	return &Logger{
		Logger: logger,
	}
}
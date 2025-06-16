package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const loggerKey = "logger"

var defaultLevel = zap.NewAtomicLevelAt(zap.InfoLevel)

func Load() *zap.SugaredLogger {
	encCfg := zapcore.EncoderConfig{
		MessageKey:       "msg",
		LevelKey:         "level",
		TimeKey:          "time",
		FunctionKey:      zapcore.OmitKey,
		LineEnding:       zapcore.DefaultLineEnding,
		ConsoleSeparator: " | ",
		EncodeLevel:      zapcore.CapitalColorLevelEncoder,
		EncodeTime:       zapcore.ISO8601TimeEncoder,
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:         "console",
		EncoderConfig:    encCfg,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger := zap.Must(cfg.Build()).Sugar()

	return logger
}

func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	if log, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return log
	}
	return zap.NewNop().Sugar()
}

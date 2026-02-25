package logger

import (
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *zap.Logger
	once     sync.Once
)

type ctxLogKey struct{}

func Init() {
	once.Do(func() {
		consoleConfig := zapcore.EncoderConfig{
			TimeKey:        "t",
			LevelKey:       "l",
			MessageKey:     "m",
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeDuration: zapcore.NanosDurationEncoder,
		}
		consoleOutput := zapcore.Lock(os.Stdout)
		consoleCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(consoleConfig),
			consoleOutput,
			zap.InfoLevel,
		)
		core := zapcore.NewTee(consoleCore)
		instance = zap.New(core, zap.AddCaller())
	})
}
func L() *zap.Logger {
	return instance
}
func ToContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLogKey{}, l)
}
func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxLogKey{}).(*zap.Logger); ok {
		return l
	}
	return L()
}

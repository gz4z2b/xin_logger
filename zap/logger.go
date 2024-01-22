package zap

import (
	"github.com/gz4z2b/xinlogger/types"
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(logger *zap.Logger) types.XinLogger {
	return &ZapLogger{
		logger: logger,
	}
}

func (l *ZapLogger) Info(msg string, args ...types.Field) {
	l.logger.Info(msg, argsToField(args)...)
}

func (l *ZapLogger) Warn(msg string, args ...types.Field) {
	l.logger.Warn(msg, argsToField(args)...)
}

func (l *ZapLogger) Error(msg string, args ...types.Field) {
	l.logger.Error(msg, argsToField(args)...)
}

func (l *ZapLogger) Debug(msg string, args ...types.Field) {
	l.logger.Debug(msg, argsToField(args)...)
}

func argsToField(args []types.Field) []zap.Field {
	result := make([]zap.Field, len(args))
	for key, arg := range args {
		result[key] = zap.Any(arg.Key, arg.Value)
	}
	return result
}

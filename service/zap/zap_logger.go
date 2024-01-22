package zap

import (
	"github.com/gz4z2b/xin_logger/service"
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger  *zap.Logger
	logPath string
}

func NewZapLogger(logPath string, logger *zap.Logger) service.XinLogger {
	return &ZapLogger{
		logger:  logger,
		logPath: logPath,
	}
}

func (l *ZapLogger) Info(msg string, args ...map[string]any) {
	l.logger.Info(msg, args...)
}

func (l *ZapLogger) Warn(msg string, args ...map[string]any) {
	panic("not implemented") // TODO: Implement
}

func (l *ZapLogger) Error(msg string, args ...map[string]any) {
	panic("not implemented") // TODO: Implement
}

func (l *ZapLogger) Debug(msg string, args ...map[string]any) {
	panic("not implemented") // TODO: Implement
}

func argsToField(args []map[string]any) zap.Field {
	result := make([]zap.Field, len(args))
	for _, arg := range args {
		for key, val := range arg {
			result = append(result, zap.Any(key, val))
		}
	}
}

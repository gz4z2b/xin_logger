package zap

import (
	"github.com/google/uuid"
	"github.com/gz4z2b/xinlogger/types"
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
	conf   types.Conf
}

func NewZapLogger(logger *zap.Logger, conf types.Conf) types.XinLogger {
	return &ZapLogger{
		logger: logger,
		conf:   conf,
	}
}

func (l *ZapLogger) Info(msg string, args ...types.Field) {
	l.logger.Info(msg, l.argsToField(args)...)
}

func (l *ZapLogger) Warn(msg string, args ...types.Field) {
	l.logger.Warn(msg, l.argsToField(args)...)
}

func (l *ZapLogger) Error(msg string, args ...types.Field) {
	l.logger.Error(msg, l.argsToField(args)...)
}

func (l *ZapLogger) Debug(msg string, args ...types.Field) {
	l.logger.Debug(msg, l.argsToField(args)...)
}

func (l *ZapLogger) FlushTraceId() error {
	if l.conf.TraceId == "" {
		traceId, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		l.conf.TraceId = traceId.String()
	}
	return nil
}

func (l *ZapLogger) argsToField(args []types.Field) []zap.Field {
	result := make([]zap.Field, len(args))
	haveTraceId := false
	for key, arg := range args {
		if arg.Key == l.conf.TraceIdKey {
			haveTraceId = true
		}
		result[key] = zap.Any(arg.Key, arg.Value)
	}
	if !haveTraceId {
		if l.conf.TraceId == "" {
			err := l.FlushTraceId()
			if err != nil {
				panic(err)
			}
		}
		result = append(result, zap.String(l.conf.TraceIdKey, l.conf.TraceId))
	}
	return result
}

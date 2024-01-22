package zap

import (
	"github.com/gz4z2b/xinlogger/types"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(conf types.Conf) types.XinLogger {
	writeSyncer := getLogWriter(conf)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, syncLogLevel(conf.EnableLevel))

	logger := zap.New(core, zap.AddCaller())
	return NewZapLogger(logger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	return encoder
}

func getLogWriter(conf types.Conf) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename: conf.LogPath,
		MaxSize:  conf.MaxSize,
		MaxAge:   conf.MaxAge,
		Compress: false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func syncLogLevel(level types.Level) zapcore.Level {
	switch level {
	case types.DebugLevel:
		return zapcore.DebugLevel
	case types.InfoLevel:
		return zapcore.InfoLevel
	case types.WarnLevel:
		return zapcore.WarnLevel
	case types.ErrorLevel:
		return zapcore.ErrorLevel
	default:
		panic("日志等级不合法")
	}
}

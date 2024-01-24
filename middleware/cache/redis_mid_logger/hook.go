package redismidlogger

import (
	"context"

	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLoggerHook struct {
	// logFun 记录日志方法，传入执行的语句以及花费时间（毫秒）
	logFun func(cmd string, milliSeconds int)
}

func NewRedisLoggerHook(logFun func(cmd string, milliSeconds int)) redis.Hook {
	return &RedisLoggerHook{
		logFun: logFun,
	}
}

func (l *RedisLoggerHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (l *RedisLoggerHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		startTime := time.Now()

		next(ctx, cmd)

		useTime := time.Since(startTime).Milliseconds()
		l.logFun(cmd.String(), int(useTime))
		return nil
	}
}

func (l *RedisLoggerHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		startTime := time.Now()

		next(ctx, cmds)

		useTime := time.Since(startTime).Milliseconds()
		var cmdStr string
		for _, cmd := range cmds {
			cmdStr += cmd.String() + "\n"
		}
		l.logFun(cmdStr, int(useTime))
		return nil
	}
}

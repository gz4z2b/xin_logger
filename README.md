# xinlogger
一个通用的go日志解决方案

# 快速入门
## 日志包使用
```go
import (
    "github.com/gz4z2b/xinlogger/types"
    "github.com/gz4z2b/xinlogger"
)

func main() {
    logger, err := xinlogger.NewLogger(types.Conf{
        // 日志路径,根目录为当前目录
		LogPath: "./test.log",
        // 日志框架,目前只支持zap
		Type:    "zap",
        // 日志级别,默认为info 
		//EnableLevel: types.DebugLevel,
        // 文件最大大小，Mb默认为5
	    //MaxSize: 10,
	    // 日志最大保存时间，天,默认为30
	    //MaxAge: 30,
	    // 最大保留文件个数,默认为10
	    //MaxBackups: 10,
        // 追踪id的初始化，例如接口有传入就在初始化时就带上
        //TraceId:"",
        // 追踪id的日志记录key值
	    TraceIdKey: "traceId", 
	})
    if err != nil {
        panic(err)
    }

    logger.Info("test", types.NewField("第一个", "第一个Info"))
    logger.Warn("test", types.NewField("第一个", "第一个Warn"))
    logger.Error("test", types.NewField("第一个", "第一个Error"))
    logger.Debug("test", types.NewField("第一个", "第一个DEBUG"))

    defer func() {
        // 刷新traceId，下一次请求进来就会是新的traceId
        logger.FlushTraceId()
    }()
}
```

## gin的http访问日志中间件
```go

import(
    "github.com/gz4z2b/xinlogger/types"
    "github.com/gz4z2b/xinlogger"
    "github.com/gz4z2b/xinlogger/middleware/http/ginmidlogger"
)

func main() {
    // 替换为自己的server初始化方法
    server := gin.Default()
    // 替换为自己的logger初始化
    logger, err := xinlogger.NewLogger(types.Conf{
        // 日志路径,根目录为当前目录
		LogPath: "./test.log",
        // 日志框架,目前只支持zap
		Type:    "zap",
	    TraceIdKey: "traceId", 
	})
    if err != nil {
        panic(err)
    }
	server.Use(ginmidlogger.NewBuilder(
        func(ctx context.Context, al *ginmidlogger.AccessLog) {
			logger.Debug("access log", types.NewField("log", al))
            // 作为http请求日志，记得最后按照自己的需求手动刷新一下traceId，目前做不到自动刷新
            logger.FlushTraceId()
        },
		}).AllowReqBody(true).AllowRespBody(true).Build()
    )
}
```

## gorm的操作日志插件

```go
    import(
        "github.com/gz4z2b/xinlogger/types"
        "github.com/gz4z2b/xinlogger"
        "github.com/gz4z2b/xinlogger/middleware/database/gormmidlogger"
    )
    func main() {
        // 替换为自己的数据库连接 begin
        connection := mysql.Open("user:password@tcp(host:port)/database?charset=charset")
	    db, err := gorm.Open(connection, &gorm.Config{})
	    if err != nil {
	    	t.Error(err)
	    }
        // 替换为自己的数据库连接 end

        // 替换为自己的logger初始化方法
	    logger, err := NewLogger(types.Conf{
	    	LogPath: "./sql.log",
	    	Type:    "zap",
	    	//EnableLevel: types.DebugLevel,
	    })
	    if err != nil {
	    	t.Error(err)
	    }
	    
        // 注册插件
	    db.Use(gormmidlogger.NewSqlLoggerMid(func(sql string, rows, seconds int) {
	    	logger.Info("数据库操作", types.NewField("sql", sql), types.NewField("effect_rows", rows), types.NewField("use_seconds", seconds))
	    }))
    }
    
```

## redis操作日志

```go
    import(
        "github.com/gz4z2b/xinlogger/types"
        "github.com/gz4z2b/xinlogger"
        "github.com/gz4z2b/xinlogger/middleware/cache/redismidlogger"
    )

    func main() {
        // 替换为自己的redis初始化
        r := redis.NewClient(&redis.Options{
	    	Addr:     host,
	    	Password: password, // no password set
	    	DB:       db,        // use default DB
	    })
        // 替换为日志的初始化
	    logger, err := NewLogger(types.Conf{
	    	LogPath: "./redis.log",
	    	Type:    "zap",
	    	MaxSize: 1,
	    	//EnableLevel: types.DebugLevel,
	    })
	    if err != nil {
	    	t.Error(err)
	    }

        // 注册插件
	    r.AddHook(redismidlogger.NewRedisLoggerHook(func(cmd string, milliSeconds int) {
	    	logger.Info("redis操作", types.NewField("命令", cmd), types.NewField("用时", milliSeconds))
	    }))
    }
```
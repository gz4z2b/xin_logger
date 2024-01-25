# xinlogger
一个通用的go日志接口

# 快速使用
## 日志包

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
	    //MaxSize: 10
	    // 日志最大保存时间，天,默认为30
	    //MaxAge: 30
	    // 最大保留文件个数,默认为10
	    //MaxBackups: 10
	},)
}
```
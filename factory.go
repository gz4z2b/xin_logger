package xinlogger

import (
	"errors"

	"github.com/gz4z2b/xinlog/types"
	"github.com/gz4z2b/xinlog/zap"
)

const ZAP_TYPE string = "zap"

var Err_TypeNotSupport = errors.New("日志框架不支持")

func NewLogger(conf types.Conf) (types.XinLogger, error) {
	err := types.CheckConf(&conf)
	if err != nil {
		return nil, err
	}
	switch conf.Type {
	case ZAP_TYPE:
		return zap.NewLogger(conf), nil
	default:
		return nil, Err_TypeNotSupport
	}
}

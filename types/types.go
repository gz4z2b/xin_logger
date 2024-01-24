package types

import (
	"errors"
)

var Err_LevelIllegal = errors.New("日志等级不合法")
var Err_PathNotNull = errors.New("日志路径不能为空")

type Conf struct {
	LogPath string
	Type    string
	// 文件最大大小，Mb
	MaxSize int
	// 日志最大保存时间，天
	MaxAge int
	// 最大保留文件个数
	MaxBackups int
	// 记录最低等级
	EnableLevel Level
	TraceId     string
	TraceIdKey  string
}

type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.

	_minLevel = DebugLevel
	_maxLevel = ErrorLevel
)

const default_maxage = 30
const default_maxsize = 5
const default_max_files = 10

func CheckConf(conf *Conf) error {
	if conf.EnableLevel < _minLevel || conf.EnableLevel > _maxLevel {
		return Err_LevelIllegal
	}
	if conf.MaxAge == 0 {
		conf.MaxAge = default_maxage
	}
	if conf.MaxSize == 0 {
		conf.MaxSize = default_maxsize
	}
	if conf.TraceIdKey == "" {
		conf.TraceIdKey = "traceId"
	}
	if conf.MaxBackups == 0 {
		conf.MaxBackups = default_max_files
	}
	if conf.LogPath == "" {
		return Err_PathNotNull
	}
	if conf.EnableLevel == 0 {
		conf.EnableLevel = InfoLevel
	}

	return nil
}

type Field struct {
	Key   string
	Value any
}

func NewField(key string, value any) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

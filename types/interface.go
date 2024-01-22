package types

type XinLogger interface {
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
	Debug(msg string, args ...Field)
}

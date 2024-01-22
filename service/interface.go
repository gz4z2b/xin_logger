package service

type XinLogger interface {
	Info(msg string, map[string]any)
	Warn(msg string, map[string]any)
	Error(msg string, map[string]any)
	Debug(msg string, map[string]any)
}

package model

type LogHandler interface {
	Error(data ...interface{})
	Warn(data ...interface{})
	Info(data ...interface{})
	Debug(data ...interface{})
}

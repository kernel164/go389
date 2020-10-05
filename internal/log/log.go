package log

import (
	"log"

	"github.com/kernel164/go389/internal/model"
	//log "gopkg.in/inconshreveable/log15.v2"
)

var l model.LogHandler

func Init(name string) error {
	logHandler, err := getLogHandler(name)
	if err != nil {
		return err
	}
	l = logHandler
	return nil
}

func getLogHandler(name string) (model.LogHandler, error) {
	return NewDefaultLogHandler(name)
}

func Error(data ...interface{}) {
	l.Error(data...)
}

func Warn(data ...interface{}) {
	l.Warn(data...)
}

func Info(data ...interface{}) {
	l.Info(data...)
}

func Debug(data ...interface{}) {
	l.Debug(data...)
}

type DefaultLogHandler struct {
	model.LogHandler
}

func NewDefaultLogHandler(name string) (model.LogHandler, error) {
	return DefaultLogHandler{}, nil
}

func (h DefaultLogHandler) Error(data ...interface{}) {
	log.Println(data...)
}

func (h DefaultLogHandler) Warn(data ...interface{}) {
	log.Println(data...)
}

func (h DefaultLogHandler) Info(data ...interface{}) {
	log.Println(data...)
}

func (h DefaultLogHandler) Debug(data ...interface{}) {
	log.Println(data...)
}

package intf

import (
	"log"
)

type LoggerProvider interface {
	Provider
	GetStdLogger() *log.Logger
	Log(keyvals ...interface{}) error
	Trace(msg string, keyvals ...interface{})
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	Fatal(msg string, keyvals ...interface{})
	Panic(msg string, keyvals ...interface{})
}

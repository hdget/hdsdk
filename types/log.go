package types

import "log"

type LogProvider interface {
	Provider

	GetStdLogger() *log.Logger

	// Log add go-kit log compatibility
	Log(keyvals ...interface{}) error

	Trace(msg string, keyvals ...interface{})
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	Fatal(msg string, keyvals ...interface{})
	Panic(msg string, keyvals ...interface{})
}

// log ability
const (
	_ SdkType = SdkCategoryLog + iota
	LibLogZerolog
)

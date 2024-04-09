// Package zerolog
// @Title  logger capability of zerolog
// @Description  zerolog implementation of logger capability
// @Author  Ryan Fan 2021-06-09
// @Update  Ryan Fan 2021-06-09
package zerolog

import (
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/hdget/hdutils"
	"github.com/rs/zerolog"
	"log"
	"strings"
)

type zerologLoggerProvider struct {
	logger zerolog.Logger
}

const (
	defaultCallerSkipFrameCount = 1 // 缺省的忽略帧数目

)

// New initialize zerolog instance
func New(conf *zerologProviderConfig) (intf.LoggerProvider, error) {
	// 设置日志级别
	switch strings.ToLower(conf.Level) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// new console logger
	consoleLogger := newConsoleLogger()

	// 设置多个输出, 输出到rotateLogs和stdout
	rotateLogger, err := newRotateLogger(conf)
	if err != nil {
		return nil, err
	}

	// 多个日志通道输出
	multi := zerolog.MultiLevelWriter(rotateLogger, consoleLogger)

	// 给zerorlogger和stdlogger实例赋值
	provider := &zerologLoggerProvider{logger: zerolog.New(multi).With().Timestamp().Logger()}

	return provider, nil
}

func (p zerologLoggerProvider) Init(args ...any) error {
	panic("implement me")
}

func (p zerologLoggerProvider) GetStdLogger() *log.Logger {
	return log.New(p.logger, "stdlog: ", log.Lshortfile|log.Ldate|log.Ltime)
}

func (p *zerologLoggerProvider) Log(keyvals ...interface{}) error {
	msgValue, errValue, fields := hdutils.ParseArgs(keyvals...)
	p.logger.Trace().Caller(defaultCallerSkipFrameCount).Err(errValue).Fields(fields).Msg(msgValue)
	return nil
}

func (p *zerologLoggerProvider) Trace(msg string, keyvals ...interface{}) {
	_, errValue, fields := hdutils.ParseArgs(keyvals...)
	p.logger.Trace().Caller(defaultCallerSkipFrameCount).Err(errValue).Fields(fields).Msg(msg)
}

func (p *zerologLoggerProvider) Debug(msg string, keyvals ...interface{}) {
	_, errValue, fields := hdutils.ParseArgs(keyvals...)
	p.logger.Debug().Caller(defaultCallerSkipFrameCount).Err(errValue).Fields(fields).Msg(msg)
}

func (p *zerologLoggerProvider) Info(msg string, keyvals ...interface{}) {
	_, errValue, fields := hdutils.ParseArgs(keyvals...)
	p.logger.Info().Caller(defaultCallerSkipFrameCount).Err(errValue).Fields(fields).Msg(msg)
}

func (p *zerologLoggerProvider) Warn(msg string, keyvals ...interface{}) {
	_, errValue, fields := hdutils.ParseArgs(keyvals...)
	p.logger.Warn().Caller(defaultCallerSkipFrameCount).Err(errValue).Fields(fields).Msg(msg)
}

func (p *zerologLoggerProvider) Error(msg string, keyvals ...interface{}) {
	_, errValue, fields := hdutils.ParseArgs(keyvals...)
	p.logger.Error().Caller(defaultCallerSkipFrameCount).Err(errValue).Fields(fields).Msg(msg)
}

func (p *zerologLoggerProvider) Fatal(msg string, keyvals ...interface{}) {
	_, errValue, fields := hdutils.ParseArgs(keyvals...)
	p.logger.Fatal().Caller(defaultCallerSkipFrameCount).Err(errValue).Fields(fields).Msg(msg)
}

func (p *zerologLoggerProvider) Panic(msg string, keyvals ...interface{}) {
	_, errValue, fields := hdutils.ParseArgs(keyvals...)
	p.logger.Panic().Caller(defaultCallerSkipFrameCount).Err(errValue).Fields(fields).Msg(msg)
}

// @Title  log capability of zerolog
// @Description  zerolog implementation of log capability
// @Author  Ryan Fan 2021-06-09
// @Update  Ryan Fan 2021-06-09
package zerolog

import (
	"fmt"
	"github.com/hdget/hdsdk/types"
	"github.com/hdget/hdsdk/utils"
	"github.com/rs/zerolog"
	stdlog "log"
	"os"
	"path/filepath"
	"strings"
)

type ZerologProvider struct{}

// 缺省的忽略帧数目
const (
	DEFAULT_CALLER_SKIP_FRAME_COUNT = 1
)

var (
	zeroLogger zerolog.Logger
	stdLogger  *stdlog.Logger
	_          types.Provider    = (*ZerologProvider)(nil)
	_          types.LogProvider = (*ZerologProvider)(nil)
	ErrKeyName                   = "err"
)

// Init	implements types.Provider interface, used to initialize the capability
// @author	Ryan Fan	(2021-06-09)
// @param	baseconf.Configer	root config interface to extract config info
// @return	error
func (c *ZerologProvider) Init(rootConfiger types.Configer, logger types.LogProvider, args ...interface{}) error {
	// 获取日志配置信息
	config, err := getLogConfig(rootConfiger)
	if err != nil {
		return err
	}

	// 设置多个输出, 输出到rotateLogs和stdout
	rotateLogs, err := newRotateLogs(config.Rotate, config.Filename)
	if err != nil {
		return err
	}

	// 设置日志级别
	switch strings.ToLower(config.Level) {
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

	// 格式化
	writer := getCustomizedConsoleWriter()

	// 多个日志通道输出
	multi := zerolog.MultiLevelWriter(rotateLogs, writer)

	// 给zerorlogger和stdlogger实例赋值
	zeroLogger = zerolog.New(multi).With().Timestamp().Logger()
	stdLogger = stdlog.New(zeroLogger, "stdlog: ", stdlog.Lshortfile|stdlog.Ldate|stdlog.Ltime)
	return nil
}

func (c *ZerologProvider) GetStdLogger() *stdlog.Logger {
	return stdLogger
}

func (c *ZerologProvider) Log(keyvals ...interface{}) error {
	msgValue, errValue, fields := utils.ParseArgsWithMsgError(keyvals...)
	zeroLogger.Trace().Caller(DEFAULT_CALLER_SKIP_FRAME_COUNT).Err(errValue).Fields(fields).Msg(msgValue)
	return nil
}

func (c *ZerologProvider) Trace(msg string, keyvals ...interface{}) {
	errValue, fields := utils.ParseArgsWithError(keyvals...)
	zeroLogger.Trace().Caller(DEFAULT_CALLER_SKIP_FRAME_COUNT).Err(errValue).Fields(fields).Msg(msg)
}

func (c *ZerologProvider) Debug(msg string, keyvals ...interface{}) {
	errValue, fields := utils.ParseArgsWithError(keyvals...)
	zeroLogger.Debug().Caller(DEFAULT_CALLER_SKIP_FRAME_COUNT).Err(errValue).Fields(fields).Msg(msg)
}

func (c *ZerologProvider) Info(msg string, keyvals ...interface{}) {
	errValue, fields := utils.ParseArgsWithError(keyvals...)
	zeroLogger.Info().Caller(DEFAULT_CALLER_SKIP_FRAME_COUNT).Err(errValue).Fields(fields).Msg(msg)
}

func (c *ZerologProvider) Warn(msg string, keyvals ...interface{}) {
	errValue, fields := utils.ParseArgsWithError(keyvals...)
	zeroLogger.Warn().Caller(DEFAULT_CALLER_SKIP_FRAME_COUNT).Err(errValue).Fields(fields).Msg(msg)
}

func (c *ZerologProvider) Error(msg string, keyvals ...interface{}) {
	errValue, fields := utils.ParseArgsWithError(keyvals...)
	zeroLogger.Error().Caller(DEFAULT_CALLER_SKIP_FRAME_COUNT).Err(errValue).Fields(fields).Msg(msg)
}

func (c *ZerologProvider) Fatal(msg string, keyvals ...interface{}) {
	errValue, fields := utils.ParseArgsWithError(keyvals...)
	zeroLogger.Fatal().Caller(DEFAULT_CALLER_SKIP_FRAME_COUNT).Err(errValue).Fields(fields).Msg(msg)
}

func (c *ZerologProvider) Panic(msg string, keyvals ...interface{}) {
	errValue, fields := utils.ParseArgsWithError(keyvals...)
	zeroLogger.Panic().Caller(DEFAULT_CALLER_SKIP_FRAME_COUNT).Err(errValue).Fields(fields).Msg(msg)
}

// 自定义输出格式
func getCustomizedConsoleWriter() zerolog.ConsoleWriter {
	// 标准输出格式
	w := zerolog.ConsoleWriter{
		Out:     os.Stdout,
		NoColor: true,
		// TimeFormat: time.RFC3339,
		TimeFormat: "2006/01/02 15:04:05",
	}

	// format functions
	w.FormatMessage = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return fmt.Sprintf("msg=\"%s\"", i)
	}

	w.FormatCaller = func(i interface{}) string {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			if cwd, err := os.Getwd(); err == nil {
				if rel, err := filepath.Rel(cwd, c); err == nil {
					c = rel
				}
			}
		}
		return c
	}

	return w
}

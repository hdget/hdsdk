package hdutils

import (
	"fmt"
	"log"
	"strings"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "DBG"
	LogLevelWarn  LogLevel = "WRN"
	LogLevelError LogLevel = "ERR"
	LogLevelFatal LogLevel = "FTL"
)

func LogDebug(msg string, keyvals ...interface{}) {
	logPrint(LogLevelDebug, msg, keyvals...)
}

func LogWarn(msg string, keyvals ...interface{}) {
	logPrint(LogLevelWarn, msg, keyvals...)
}

func LogError(msg string, keyvals ...interface{}) {
	logPrint(LogLevelError, msg, keyvals...)
}

func LogFatal(msg string, keyvals ...interface{}) {
	logPrint(LogLevelFatal, msg, keyvals...)
}

// ParseArgs  解析error和message用统一格式展示出来
func ParseArgs(keyvals ...interface{}) (string, error, map[string]interface{}) {
	countArgs := len(keyvals)
	// 如果可变参数个数为0，肯定没有error
	if countArgs == 0 {
		return "", nil, nil
	}

	var errValue error
	var msgValue string
	args := make(map[string]interface{})
	for i := 0; i < countArgs-1; i = i + 2 {
		// 第i个值作为map的key, 第i+1个值作为map的value
		k, ok := keyvals[i].(string)
		if !ok {
			continue
		}

		switch strings.ToLower(k) {
		case "level", "caller":
			// do nothing
		case "err", "error": // err and error must be error type
			switch v := keyvals[i+1].(type) {
			case error:
				errValue = v
			default:
				errValue = fmt.Errorf("%v", keyvals[i+1])
			}
		case "msg", "message":
			msgValue = fmt.Sprintf("%v", keyvals[i+1])
		default:
			args[k] = keyvals[i+1]
		}
	}
	return msgValue, errValue, args
}

// logPrint log structure message and key values
func logPrint(level LogLevel, msg string, keyvals ...interface{}) {
	_, errValue, fields := ParseArgs(keyvals...)

	outputs := make([]string, 0)
	for k, v := range fields {
		outputs = append(outputs, fmt.Sprintf("%s=\"%v\"", k, v))
	}

	logFn := log.Printf
	if level == LogLevelFatal {
		logFn = log.Fatalf
	}

	if len(outputs) > 0 {
		if errValue != nil {
			logFn("%s msg=\"%s\" %s error=\"%v\"", level, msg, strings.Join(outputs, " "), errValue)
		} else {
			logFn("%s msg=\"%s\" %s", level, msg, strings.Join(outputs, " "))
		}
	} else {
		if errValue != nil {
			logFn("%s msg=\"%s\" error=\"%v\"", level, msg, errValue)
		} else {
			logFn("%s msg=\"%s\"", level, msg)
		}
	}
}

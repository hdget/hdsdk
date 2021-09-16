package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"strings"
)

// Fatal 结构化的Fatal
func Fatal(msg string, keyvals ...interface{}) {
	errValue, fields := ParseArgsWithError(keyvals...)

	outputs := make([]string, 0)
	for k, v := range fields {
		outputs = append(outputs, fmt.Sprintf("%s=\"%v\"", k, v))
	}

	if len(outputs) > 0 {
		if errValue != nil {
			log.Fatalf("FTL msg=\"%s\" %s error=\"%v\"", msg, strings.Join(outputs, " "), errValue)
		} else {
			log.Fatalf("FTL msg=\"%s\" %s", msg, strings.Join(outputs, " "))
		}
	} else {
		if errValue != nil {
			log.Fatalf("FTL msg=\"%s\" error=\"%v\"", msg, errValue)
		} else {
			log.Fatalf("FTL msg=\"%s\"", msg)
		}
	}
}

// Print 结构化的log
func Print(level, msg string, keyvals ...interface{}) {
	errValue, fields := ParseArgsWithError(keyvals...)

	outputs := make([]string, 0)
	for k, v := range fields {
		outputs = append(outputs, fmt.Sprintf("%s=\"%v\"", k, v))
	}

	if len(outputs) > 0 {
		if errValue != nil {
			log.Printf("%s msg=\"%s\" %s error=\"%v\"", level, msg, strings.Join(outputs, " "), errValue)
		} else {
			log.Printf("%s msg=\"%s\" %s", level, msg, strings.Join(outputs, " "))
		}
	} else {
		if errValue != nil {
			log.Printf("%s msg=\"%s\" error=\"%v\"", level, msg, errValue)
		} else {
			log.Printf("%s msg=\"%s\"", level, msg)
		}
	}
}

// ParseArgsWithError  将可变参数转换成map, 其中有err关键字返回error
// @author Ryan Fan
// @param  variadic arguments, key/value键值对，可变参数个数必须为偶数
// @return 错误值
// @return 除错误值以外的其他参数key->value对
func ParseArgsWithError(keyvals ...interface{}) (error, map[string]interface{}) {
	countArgs := len(keyvals)
	// 如果可变参数个数为0，肯定没有error
	if countArgs == 0 {
		return nil, nil
	}

	var errValue error
	args := make(map[string]interface{})
	for i := 0; i < countArgs; i = i + 2 {
		// 如果下一个值的index小于参数个数，继续进行判断
		if i+1 < countArgs {
			// 第i个值作为map的key, 第i+1个值作为map的value
			k, ok := keyvals[i].(string)
			if !ok {
				continue
			}

			switch k {
			case "err", "error":
				v, ok := keyvals[i+1].(error)
				if ok {
					errValue = v
				} else {
					// if next value is not an error, try convert it's string representation to error
					errValue = errors.New(fmt.Sprintf("%v", keyvals[i+1]))
				}
			default:
				args[k] = keyvals[i+1]
			}
		}
	}

	return errValue, args
}

// ParseArgsWithMsgError  将可变参数转换成map, 其中有err关键字返回error, 有msg关键子返回msg value
// @author Ryan Fan
// @param  variadic arguments, key/value键值对，可变参数个数必须为偶数
// @return 错误值
// @return 除错误值以外的其他参数key->value对
func ParseArgsWithMsgError(keyvals ...interface{}) (string, error, map[string]interface{}) {
	countArgs := len(keyvals)
	// 如果可变参数个数为0，肯定没有error
	if countArgs == 0 {
		return "", nil, nil
	}

	var msgValue string
	var errValue error
	args := make(map[string]interface{})
	for i := 0; i < countArgs; i = i + 2 {
		// 如果下一个值的index小于参数个数，继续进行判断
		if i+1 < countArgs {
			// 第i个值作为map的key, 第i+1个值作为map的value
			k, ok := keyvals[i].(string)
			if !ok {
				continue
			}

			switch k {
			case "msg":
				if v, ok := keyvals[i+1].(string); ok {
					msgValue = v
				}
			case "err":
				if v, ok := keyvals[i+1].(error); ok {
					errValue = v
				}
			default:
				args[k] = keyvals[i+1]
			}
		}
	}
	return msgValue, errValue, args
}

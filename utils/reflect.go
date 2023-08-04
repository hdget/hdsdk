package utils

import (
	"reflect"
	"runtime"
	"strings"
)

// GetFuncName 从函数实例获取函数名
func GetFuncName(fn any) string {
	tokens := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
	return strings.Split(tokens[len(tokens)-1], "-")[0]
}

func GetStructName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

// GetVarName 获取变量的名字
func GetVarName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

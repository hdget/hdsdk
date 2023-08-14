package utils

import (
	"fmt"
	"github.com/pkg/errors"
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

// StructSet 将结构中的字段设置为某个值
func StructSet(obj any, typ any, val any) error {
	foundField := false
	numField := reflect.TypeOf(obj).Elem().NumField()
	for i := 0; i < numField; i++ {
		field := reflect.ValueOf(obj).Field(i)
		if field.Type() == reflect.TypeOf(typ) {
			if !field.CanSet() {
				return errors.New("can not set")
			}
			foundField = true
			field.Set(reflect.ValueOf(val))
		}
	}

	if !foundField {
		return fmt.Errorf("module need inherits %#v", reflect.TypeOf(typ))
	}
	return nil
}

func StructGetReceivers(obj any, fn any) map[string]any {
	receivers := make(map[string]any)

	// common.ServiceInvocationHandler(nil)
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	numMethod := v.NumField()
	for i := 0; i < numMethod; i++ {
		tt := t.Method(i)
		vv := v.Method(i)
		if vv.Type().ConvertibleTo(reflect.TypeOf(fn)) {
			receivers[tt.Name] = vv.Interface()
		}
	}

	return receivers
}

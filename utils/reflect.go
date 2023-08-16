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

func GetStructName(obj any) string {
	if obj == nil {
		return ""
	}

	var st reflect.Type
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		st = t.Elem()
	} else {
		st = t
	}
	if st.Kind() != reflect.Struct {
		return ""
	}
	return st.Name()
}

// GetVarName 获取变量的名字
func GetVarName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

// StructSetComplexField 将结构中的接口或者结构指针设置为某个值
func StructSetComplexField(obj any, emptyFieldObj any, val any) error {
	if obj == nil {
		return errors.New("nil struct")
	}

	// struct有可能是指针
	var st reflect.Type
	var sv reflect.Value
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		st = reflect.TypeOf(obj).Elem()
		sv = reflect.ValueOf(obj).Elem()
	} else {
		st = reflect.TypeOf(obj)
		sv = reflect.ValueOf(obj)
	}

	numField := st.NumField()
	for i := 0; i < numField; i++ {
		field := sv.Field(i)
		fieldType := field.Type().String()
		emptyFieldType := reflect.TypeOf(emptyFieldObj).String()
		// 找到第一个匹配类型的field设置进去
		if fieldType == emptyFieldType || "*"+fieldType == emptyFieldType {
			if !field.CanSet() {
				return errors.New("field can not set")
			}
			field.Set(reflect.ValueOf(val))
			return nil
		}
	}
	return fmt.Errorf("no field match %#v", reflect.TypeOf(emptyFieldObj))
}

func StructGetReceiverMethodsByType(receiver any, fn any) map[string]any {
	if receiver == nil {
		return nil
	}

	st := reflect.TypeOf(receiver)
	sv := reflect.ValueOf(receiver)
	numMethod := sv.NumMethod()

	receivers := make(map[string]any)
	for i := 0; i < numMethod; i++ {
		vv := sv.Method(i)
		if vv.Type().ConvertibleTo(reflect.TypeOf(fn)) {
			receivers[st.Method(i).Name] = vv.Interface()
		}
	}
	return receivers
}

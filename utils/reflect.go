package utils

import (
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

func StructSetComplexField(obj any, nilField any, val any) error {
	if obj == nil {
		return errors.New("nil struct")
	}

	if val == nil {
		return errors.New("cannot set to zero value")
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
		emptyFieldType := reflect.TypeOf(nilField).String()
		// 找到第一个匹配类型的field设置进去
		if fieldType == emptyFieldType || "*"+fieldType == emptyFieldType {
			if !field.CanSet() {
				return errors.New("field cannot set")
			}

			field.Set(reflect.ValueOf(val))
			return nil
		}
	}
	return errors.New("filed not found")
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

// GetFuncSignature 获取函数签名信息
func GetFuncSignature(fn any) string {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		return ""
	}

	buf := strings.Builder{}
	buf.WriteString("func (")
	for i := 0; i < t.NumIn(); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(t.In(i).String())
	}
	buf.WriteString(")")
	if numOut := t.NumOut(); numOut > 0 {
		if numOut > 1 {
			buf.WriteString(" (")
		} else {
			buf.WriteString(" ")
		}
		for i := 0; i < t.NumOut(); i++ {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(t.Out(i).String())
		}
		if numOut > 1 {
			buf.WriteString(")")
		}
	}

	return buf.String()
}

package utils

import (
	"fmt"
	"testing"
)

type TestInterface interface {
	Get()
}
type testStruct struct {
	TestInterface
}

func (*testStruct) Get() {}

func TestStructSetInterfaceField(t *testing.T) {

	type args struct {
		obj       any
		filedType any
		val       any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestStructSetField",
			args: args{
				obj:       &testStruct{},
				filedType: (*TestInterface)(nil),
				val:       nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StructSetComplexField(tt.args.obj, tt.args.filedType, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("StructSetComplexField() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStructSetStructField(t *testing.T) {
	type anotherStruct struct {
		Name string
	}
	type testStruct struct {
		Another *anotherStruct
	}
	type args struct {
		obj       any
		filedType any
		val       any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestStructSetStructField",
			args: args{
				obj:       &testStruct{},
				filedType: &anotherStruct{},
				val:       &anotherStruct{Name: "xxx"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StructSetComplexField(tt.args.obj, tt.args.filedType, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("StructSetComplexField() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (*testStruct) Aaa(arg string) error {
	println(arg)
	return nil
}

type testfunc func(string) error

func TestStructGetReceiverMethods(t *testing.T) {
	type args struct {
		obj any
		fn  any
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "TestStructGetReceiverMethods",
			args: args{
				obj: &testStruct{},
				fn:  testfunc(nil),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gots := StructGetReceiverMethodsByType(tt.args.obj, tt.args.fn)
			fmt.Println(len(gots))
		})
	}
}

func TestGetFuncSignature(t *testing.T) {
	type anyFn func(interface{}) any
	type anyFn1 func(any) interface{}
	type args struct {
		fn any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				fn: anyFn(nil),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFuncSignature(anyFn(nil))
			got1 := GetFuncSignature(anyFn1(nil))
			if got != got1 {
				t.Errorf("GetFuncSignature() not equal, got: %v, got1: %v", got, got1)
			}
		})
	}
}

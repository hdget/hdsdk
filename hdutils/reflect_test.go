package hdutils

import (
	"fmt"
	"reflect"
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
			if err := Reflect().StructSet(tt.args.obj, tt.args.filedType, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("StructSet() error = %v, wantErr %v", err, tt.wantErr)
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
			if err := Reflect().StructSet(tt.args.obj, tt.args.filedType, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("StructSet() error = %v, wantErr %v", err, tt.wantErr)
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
			gots := Reflect().MatchReceiverMethods(tt.args.obj, tt.args.fn)
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
			got := Reflect().GetFuncSignature(anyFn(nil))
			got1 := Reflect().GetFuncSignature(anyFn1(nil))
			if got != got1 {
				t.Errorf("GetFuncSignature() not equal, got: %v, got1: %v", got, got1)
			}
		})
	}
}

func Test_hdReflector_InspectValue(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want *ValueMeta
	}{
		{
			name: "Test_hdReflector_InspectValue",
			args: args{
				v: nil,
			},
			want: nil,
		},
		{
			name: "Test_hdReflector_InspectValue",
			args: args{
				v: 1,
			},
			want: nil,
		},
		{
			name: "Test_hdReflector_InspectValue",
			args: args{
				v: int32(1),
			},
			want: nil,
		},
		{
			name: "Test_hdReflector_InspectValue",
			args: args{
				v: float64(0),
			},
			want: nil,
		},
		{
			name: "Test_hdReflector_InspectValue",
			args: args{
				v: "",
			},
			want: nil,
		},
		{
			name: "Test_hdReflector_InspectValue",
			args: args{
				v: "aaa",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hdReflector{}
			if got := h.InspectValue(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InspectValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

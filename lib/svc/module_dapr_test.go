package svc

import (
	"reflect"
	"testing"
)

func TestDaprModule_GetRoutes(t *testing.T) {
	type fields struct {
		BaseModule *BaseModule
	}
	type args struct {
		srcPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Route
		wantErr bool
	}{
		{
			name: "TestDaprModule_GetRoutes",
			fields: fields{
				BaseModule: NewBaseModule("testapp", "v2_aliyunoss", 1),
			},
			args: args{
				srcPath: "D:\\Codes\\workspace\\base\\service",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DaprModule{
				BaseModule: tt.fields.BaseModule,
			}
			got, err := m.GetRoutes(tt.args.srcPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRoutes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRoutes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

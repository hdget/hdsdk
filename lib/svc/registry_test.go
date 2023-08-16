package svc

import (
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/utils"
	"os"
	"testing"
)

type TestConfig struct {
	hdsdk.Config `mapstructure:",squash"`
}

func TestMain(m *testing.M) {
	var conf TestConfig
	v := hdsdk.NewConfig("test", "local").Load()
	if len(v.AllKeys()) > 0 {
		err := v.Unmarshal(&conf)
		if err != nil {
			utils.LogFatal("msg", "unmarshal config", "err", err)
		}
	}

	err := hdsdk.Initialize(&conf)
	if err != nil {
		utils.LogFatal("sdk initialize", "err", err)
	}

	m.Run()

	os.Exit(0)
}

func TestRegisterAsDaprModule(t *testing.T) {
	type args struct {
		m       any
		app     string
		version int
		options []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestRegisterAsDaprModule",
			args: args{
				m:       &testModule{},
				app:     "",
				version: 0,
				options: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterDaprModule(tt.args.m, tt.args.app, tt.args.version, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("RegisterAsDaprModule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

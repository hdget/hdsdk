package hdutils

import (
	"reflect"
	"testing"
)

func TestGetBetweenDays(t *testing.T) {
	type xargs struct {
		beginDate string
		args      []string
	}
	tests := []struct {
		name    string
		xargs   xargs
		want    []string
		wantErr bool
	}{
		{
			name: "test between days",
			xargs: xargs{
				beginDate: "2022-02-27",
				args:      []string{"2022-03-01"},
			},
			want:    []string{"2022-02-27", "2022-02-28", "2022-03-01"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBetweenDays(LayoutIsoDate, tt.xargs.beginDate, tt.xargs.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBetweenDays() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBetweenDays() got = %v, want %v", got, tt.want)
			}
		})
	}
}

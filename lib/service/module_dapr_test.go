package service

import "testing"

func Test_matchHandlerSuffix(t *testing.T) {
	type args struct {
		methodName string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "Test_matchHandlerSuffix",
			args: args{
				methodName: "isCheckedHandlerHANDler",
			},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := matchHandlerSuffix(tt.args.methodName)
			if got != tt.want {
				t.Errorf("matchHandlerSuffix() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("matchHandlerSuffix() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

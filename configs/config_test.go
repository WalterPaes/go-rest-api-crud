package configs

import (
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	type args struct {
		filenames []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Configs
		wantErr bool
	}{
		{
			name: "Should load configs without errors",
			args: args{
				filenames: []string{
					"../.env.example",
				},
			},
			want: &Configs{
				ApiPort: ":8080",
			},
			wantErr: false,
		},
		{
			name: "Should has an error when try load env filename",
			args: args{
				filenames: []string{
					"invalid",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.filenames...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

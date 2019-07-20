package timeular

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	type args struct {
		APIKey    string
		APISecret string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "prod",
			args: args{
				APISecret: "test",
				APIKey: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.args.APIKey, tt.args.APISecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

package timeular

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_trackingService_ShowCurrentTracking(t *testing.T) {
	type fields struct {
		c *Client
	}

	c, err := NewClient("", "")
	if !assert.NoError(t, err) {
		return
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "main",
			fields: fields{
				c: c,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &trackingService{
				c: tt.fields.c,
			}
			_, err := s.ShowCurrentTracking()
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err)
				return
			}
		})
	}
}

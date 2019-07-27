package timeular

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func init () {
	_ = godotenv.Load()
}

func Test_trackingService_ShowCurrentTracking(t *testing.T) {
	type fields struct {
		c *Client
	}

	c, err := NewClient(os.Getenv("TIMEULAR_API_KEY"), os.Getenv("TIMEULAR_API_SECRET"))
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
			got, err := s.ShowCurrentTracking()
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err)
				return
			}

			assert.NotNil(t, got.CurrentTracking)
		})
	}
}

func Test_trackingService_getTimeEntriesByRange(t *testing.T) {
	type fields struct {
		c *Client
	}

	c, err := NewClient(os.Getenv("TIMEULAR_API_KEY"), os.Getenv("TIMEULAR_API_SECRET"))
	if !assert.NoError(t, err) {
		return
	}

	type args struct {
		startedBefore time.Time
		stoppedAfter  time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "main",
			fields: fields{
				c: c,
			},
			args: args{
				startedBefore: time.Now(),
				stoppedAfter:  time.Now().Add(-time.Hour * 72),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &trackingService{
				c: tt.fields.c,
			}
			got, err := s.getTimeEntriesByRange(context.Background(), tt.args.startedBefore, tt.args.stoppedAfter)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTimeEntriesByRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.NotZero(t, got.TimeEntries)
		})
	}
}
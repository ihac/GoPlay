package prometheus

import (
	"reflect"
	"testing"
	"time"
)

var (
	client *PromClient
)

const (
	prometheus = "http://10.100.100.172:9090"
	metric     = "promhttp_metric_handler_requests_total"
	interval   = 1 * time.Minute
	step       = 15 * time.Second
)

func init() {
	client = &PromClient{
		Address:  prometheus,
		Metric:   metric,
		Interval: interval,
		Step:     step,
	}
}
func TestPromClient_GetAbnormalInstance(t *testing.T) {
	type args struct {
		checkpoint time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Test1",
			args: args{
				checkpoint: time.Now().Add(-9 * time.Hour),
			},
			want:    []string{"localhost:9090"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetAbnormalInstance(tt.args.checkpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("PromClient.GetAbnormalInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PromClient.GetAbnormalInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

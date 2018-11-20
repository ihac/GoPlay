package prometheus

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

// PromClient defines the client interface to query a metric over a period of time.
type PromClient struct {
	Address  string
	Metric   string
	Interval time.Duration
	Step     time.Duration
}

func (pc *PromClient) queryRange(end time.Time) (model.Value, error) {
	c, err := api.NewClient(api.Config{
		Address: pc.Address,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating api client for querying metrics: %v", err)
	}
	api := v1.NewAPI(c)

	start := end.Add(-1 * pc.Interval)
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  pc.Step,
	}
	v, err := api.QueryRange(context.Background(), pc.Metric, r)
	if err != nil {
		return nil, fmt.Errorf("error querying data over a period of time from Prometheus server: %v", err)
	}
	return v, nil
}

// GetAbnormalInstance get all abnormal instances by parsing metrics from Prometheus server.
func (pc *PromClient) GetAbnormalInstance(checkpoint time.Time) ([]string, error) {
	resp, err := pc.queryRange(checkpoint)
	if err != nil {
		return nil, err
	}
	m, ok := resp.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("error convert response of Prometheus to Matrix type")
	}

	var hosts []string
	visited := make(map[string]bool)
	for _, stream := range m {
		// maybe we should validate whether field 'instance' exists.
		instance := string(stream.Metric["instance"])
		logrus.Infof("instance %q met", instance)
		if _, isVisited := visited[instance]; isVisited {
			continue
		}
		for _, pair := range stream.Values {
			if pair.Value > 0 {
				visited[instance] = true
				hosts = append(hosts, instance)
				logrus.Warnf("instance %q identified as an abnormal host", instance)
				break
			}
		}
	}
	return hosts, nil
}

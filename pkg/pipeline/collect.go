package pipeline

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cncf-tags/green-reviews-tooling/pkg/cmd"
	"github.com/cncf-tags/green-reviews-tooling/pkg/sci"
)

// collect queries Prometheus for energy and project metrics after a benchmark
// run and computes the SCI score. Errors on individual non-energy metrics are
// non-fatal so that a missing metric doesn't discard the whole result.
func (p *Pipeline) collect(ctx context.Context, cncfProject, config, version string, durationMins int, startTime, endTime time.Time) (*sci.BenchmarkResult, error) {
	subquery := fmt.Sprintf("%dm:1m", durationMins)
	window := fmt.Sprintf("%dm", durationMins)

	keplerJoules, err := p.queryPrometheus(ctx, fmt.Sprintf(
		`sum(increase(kepler_container_joules_total{container_name="%s"}[%s]))`,
		cncfProject, subquery,
	))
	if err != nil {
		return nil, fmt.Errorf("collect: kepler_container_joules_total: %w", err)
	}

	cpuUsage, err := p.queryPrometheus(ctx, fmt.Sprintf(
		`sum(increase(container_cpu_usage_seconds_total{container="%s"}[%s]))`,
		cncfProject, subquery,
	))
	if err != nil {
		log.Printf("collect: container_cpu_usage_seconds_total: %v", err)
	}

	memRSS, err := p.queryPrometheus(ctx, fmt.Sprintf(
		`avg_over_time(container_memory_rss{container="%s"}[%s])`,
		cncfProject, window,
	))
	if err != nil {
		log.Printf("collect: container_memory_rss: %v", err)
	}

	memWS, err := p.queryPrometheus(ctx, fmt.Sprintf(
		`avg_over_time(container_memory_working_set_bytes{container="%s"}[%s])`,
		cncfProject, window,
	))
	if err != nil {
		log.Printf("collect: container_memory_working_set_bytes: %v", err)
	}

	return &sci.BenchmarkResult{
		CNCFProject: sci.ProjectInfo{
			Name:    cncfProject,
			Config:  config,
			Version: version,
		},
		Results: sci.Results{
			Metrics: []sci.Metric{
				{Name: "kepler_container_joules_total", Value: keplerJoules},
				{Name: "container_cpu_usage_seconds_total", Value: cpuUsage},
				{Name: "container_memory_rss", Value: memRSS},
				{Name: "container_memory_working_set_bytes", Value: memWS},
			},
			SCI: sci.ComputeSCI(keplerJoules, sci.DefaultSCIConfig()),
		},
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

// queryPrometheus runs a PromQL instant query via kubectl's API server proxy
// and returns the scalar result value.
func (p *Pipeline) queryPrometheus(ctx context.Context, query string) (float64, error) {
	stdout, err := p.exec(ctx, cmd.QueryPrometheus(query))
	if err != nil {
		return 0, fmt.Errorf("kubectl get --raw: %w", err)
	}
	return sci.ParsePrometheusResponse(stdout)
}

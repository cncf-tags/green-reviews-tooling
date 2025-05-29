package pipeline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/prometheus/common/model"
)

type MetericsCollectorResult struct {
	Query      string
	QueryType  string
	QueryValue string
}

type BenchmarkingCollectorResults []MetericsCollectorResult

func (p *Pipeline) computeBenchmarkingResults(ctx context.Context, q *Query) (BenchmarkingCollectorResults, error) {

	queryMap := []struct {
		mType model.ValueType
		query string
		mVal  string
	}{
		{
			query: "rate(container_cpu_usage_seconds_total[15m])",
		},
		{
			query: "avg_over_time(container_memory_rss[15m])",
		},
		{
			query: "avg_over_time(container_memory_working_set_bytes[15m])",
		},
	}

	for idx := range queryMap {
		d, warns, qErr := q.WithTimeRange(ctx, queryMap[idx].query, 15)
		if qErr != nil {
			return nil, qErr
		}

		if len(warns) > 0 {
			p.echo(ctx, fmt.Sprintf("Warnings received during query execution: %v", warns))
		}

		queryMap[idx].mType = d.Type()
		queryMap[idx].mVal = d.String()
	}

	p.echo(ctx, "Benchmarking results:")

	res := make(BenchmarkingCollectorResults, len(queryMap))

	for idx := range queryMap {
		res = append(res, MetericsCollectorResult{
			Query:      queryMap[idx].query,
			QueryType:  queryMap[idx].mType.String(),
			QueryValue: queryMap[idx].mVal,
		})
	}

	if b, err := json.MarshalIndent(res, "", "  "); err == nil {
		p.echo(ctx, string(b))
	} else {
		p.echo(ctx, fmt.Sprintf("Failed to serialize metricsCollectorRes: %v", err))
	}

	return res, nil
}

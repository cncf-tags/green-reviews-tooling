package pipeline

import (
	"context"
	"fmt"

	"github.com/prometheus/common/model"
)

type MetericsCollectorResult struct {
	Query      string
	QueryType  string
	QueryValue string
}

type BenchmarkingCollectorResults []MetericsCollectorResult

func (p *Pipeline) computeBenchmarkingResults(ctx context.Context, q *Query, benchmarkJobDurationMins int) (BenchmarkingCollectorResults, error) {

	queryMap := []struct {
		mType model.ValueType
		query string
		mVal  string
	}{
		{
			query: fmt.Sprintf("rate(container_cpu_usage_seconds_total[%dm])", benchmarkJobDurationMins),
		},
		{
			query: fmt.Sprintf("avg_over_time(container_memory_rss[%dm])", benchmarkJobDurationMins),
		},
		{
			query: fmt.Sprintf("avg_over_time(container_memory_working_set_bytes[%dm])", benchmarkJobDurationMins),
		},
	}

	for idx := range queryMap {
		d, warns, qErr := q.WithTimeRange(ctx, queryMap[idx].query, benchmarkJobDurationMins+1)
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

	for _, r := range res {
		p.echo(ctx, fmt.Sprintf("Query: %s, Type: %s, Value: %s", r.Query, r.QueryType, r.QueryValue))
	}

	return res, nil
}

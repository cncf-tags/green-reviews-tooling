package monitoring

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/prometheus/common/model"
)

func ComputeBenchmarkingResults(ctx context.Context) error {
	q, err := NewQuery(
		WithClientTimeout(10*time.Second),
		WithPrometheusAddress("http://kube-prometheus-stack-prometheus.monitoring:9090/"),
	)
	if err != nil {
		return err
	}

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
		d, qErr := q.WithTimeRange(ctx, queryMap[idx].query, 15)
		if qErr != nil {
			return qErr
		}

		queryMap[idx].mType = d.Type()
		queryMap[idx].mVal = d.String()
	}

	for idx := range queryMap {
		log.Println(strings.Repeat("-", 80))
		log.Println("Query: ", queryMap[idx].query)
		log.Printf("Type: %v\nVal: %s\n", queryMap[idx].mType, queryMap[idx].mVal)
		log.Println(strings.Repeat("-", 80))
	}

	return nil
}

package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/prometheus/common/model"
)

// Structured data types for different metric results
type MetricDataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

type MetricSeries struct {
	Labels     map[string]string `json:"labels"`
	DataPoints []MetricDataPoint `json:"data_points"`
}

type MetricsCollectorResult struct {
	Query       string         `json:"query"`
	QueryType   string         `json:"query_type"`
	MetricName  string         `json:"metric_name"`
	Series      []MetricSeries `json:"series"`
	ScalarValue *float64       `json:"scalar_value,omitempty"` // For scalar results
}

type BenchmarkingCollectorResults []MetricsCollectorResult

func (b BenchmarkingCollectorResults) WriteJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(b)
}

func extractStructuredData(value model.Value, query string) MetricsCollectorResult {
	result := MetricsCollectorResult{
		Query:     query,
		QueryType: value.Type().String(),
	}

	switch v := value.(type) {
	case model.Matrix:
		result.Series = make([]MetricSeries, len(v))
		for i, sample := range v {
			series := MetricSeries{
				Labels:     make(map[string]string),
				DataPoints: make([]MetricDataPoint, len(sample.Values)),
			}

			for label, value := range sample.Metric {
				series.Labels[string(label)] = string(value)
			}

			if metricName, exists := sample.Metric[model.MetricNameLabel]; exists {
				result.MetricName = string(metricName)
			}

			for j, pair := range sample.Values {
				series.DataPoints[j] = MetricDataPoint{
					Timestamp: pair.Timestamp.Time(),
					Value:     float64(pair.Value),
				}
			}

			result.Series[i] = series
		}

	case model.Vector:
		result.Series = make([]MetricSeries, len(v))
		for i, sample := range v {
			series := MetricSeries{
				Labels: make(map[string]string),
				DataPoints: []MetricDataPoint{{
					Timestamp: sample.Timestamp.Time(),
					Value:     float64(sample.Value),
				}},
			}

			for label, value := range sample.Metric {
				series.Labels[string(label)] = string(value)
			}

			if metricName, exists := sample.Metric[model.MetricNameLabel]; exists {
				result.MetricName = string(metricName)
			}

			result.Series[i] = series
		}

	case *model.Scalar:
		scalarVal := float64(v.Value)
		result.ScalarValue = &scalarVal
		result.Series = []MetricSeries{{
			Labels: map[string]string{},
			DataPoints: []MetricDataPoint{{
				Timestamp: v.Timestamp.Time(),
				Value:     float64(v.Value),
			}},
		}}
	}

	return result
}

func (p *Pipeline) computeBenchmarkingResults(
	ctx context.Context,
	q *Query,
	benchmarkJobDurationMins int,
	benchmarkNamespace string,
) (BenchmarkingCollectorResults, error) {
	queries := []string{
		fmt.Sprintf(
			`rate(container_cpu_usage_seconds_total{namespace="%s"}[%dm])`,
			benchmarkNamespace,
			benchmarkJobDurationMins,
		),
		fmt.Sprintf(
			`avg_over_time(container_memory_rss{namespace="%s"}[%dm])`,
			benchmarkNamespace,
			benchmarkJobDurationMins,
		),
		fmt.Sprintf(
			`avg_over_time(container_memory_working_set_bytes{namespace="%s"}[%dm])`,
			benchmarkNamespace,
			benchmarkJobDurationMins,
		),
	}

	results := make(BenchmarkingCollectorResults, 0, len(queries))

	for _, query := range queries {
		d, warns, qErr := q.WithTimeRange(ctx, query, benchmarkJobDurationMins+1)
		if qErr != nil {
			return nil, fmt.Errorf("failed to execute query '%s': %w", query, qErr)
		}

		if len(warns) > 0 {
			p.echo(ctx, fmt.Sprintf("Warnings received during query execution: %v", warns))
		}

		structuredResult := extractStructuredData(d, query)
		results = append(results, structuredResult)
	}

	return results, nil
}

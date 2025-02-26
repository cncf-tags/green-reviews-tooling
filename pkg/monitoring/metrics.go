package monitoring

import (
	"context"
	"log"
	"strings"
	"time"
)

func FetchMetrics(ctx context.Context) error {
	q, err := NewQuery(WithClientTimeout(10 * time.Second))
	if err != nil {
		return err
	}

	d, qErr := q.WithRange(ctx, "container_cpu_usage_seconds_total", 15)
	if qErr != nil {
		return qErr
	}

	log.Println(strings.Repeat("-", 80))
	log.Printf("Type: %v\nVal: %s\n", d.Type(), d.String())
	log.Println(strings.Repeat("-", 80))

	return nil
}

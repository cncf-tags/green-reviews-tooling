package pipeline

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	promModel "github.com/prometheus/common/model"
)

type Query struct {
	c v1.API
}

func NewQuery(promURL string) (*Query, error) {
	q := new(Query)

	client, err := api.NewClient(api.Config{
		Client: &http.Client{
			Timeout: 10 * time.Second, // Default timeout, can be overridden
		},
		Address: promURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create prometheus client: %w", err)
	}

	q.c = v1.NewAPI(client)

	return q, nil
}

func (q *Query) WithTimeRange(ctx context.Context, query string, dtInMinutes int) (promModel.Value, v1.Warnings, error) {
	_ctx, cancel := context.WithTimeout(ctx, time.Duration(dtInMinutes)*time.Minute)
	defer cancel()

	r := v1.Range{
		Start: time.Now().Add(-time.Duration(dtInMinutes) * time.Minute),
		End:   time.Now(),
		Step:  time.Minute,
	}

	result, warnings, err := q.c.QueryRange(_ctx, query, r)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query prometheus: %w", err)
	}

	if len(warnings) > 0 {
		return result, warnings, nil
	}

	return result, nil, nil
}

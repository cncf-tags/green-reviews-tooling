package monitoring

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	promModel "github.com/prometheus/common/model"
)

type Query struct {
	c             v1.API
	clientTimeout time.Duration
}

type option struct {
	prometheusAddress string
	clientTimeout     time.Duration
}

type MetricsClientOption func(*option) error

func WithPrometheusAddress(address string) func(*option) error {
	return func(o *option) error {
		o.prometheusAddress = address
		return nil
	}
}

func WithClientTimeout(timeout time.Duration) func(*option) error {
	return func(o *option) error {
		o.clientTimeout = timeout
		return nil
	}
}

// TODO: need to get the logger if dagger provides one!!

func NewQuery(opts ...MetricsClientOption) (*Query, error) {
	q := new(Query)

	o := &option{}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	addr, dur := "", time.Duration(0)
	if o.prometheusAddress != "" {
		addr = o.prometheusAddress
	} else {
		addr = "http://localhost:9090" // TODO: need to make a meaningful default for testing purposes
	}

	if o.clientTimeout != 0 {
		dur = o.clientTimeout
	} else {
		dur = time.Minute
	}

	client, err := api.NewClient(api.Config{
		Client: &http.Client{
			Timeout: dur,
		},
		Address: addr,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create prometheus client: %w", err)
	}

	q.c = v1.NewAPI(client)
	q.clientTimeout = dur

	return q, nil
}

func (q *Query) WithTimeRange(ctx context.Context, query string, dtInMinutes int) (promModel.Value, error) {
	_ctx, cancel := context.WithTimeout(ctx, q.clientTimeout)
	defer cancel()

	r := v1.Range{
		Start: time.Now().Add(-time.Duration(dtInMinutes) * time.Minute),
		End:   time.Now(),
		Step:  time.Minute,
	}

	result, warnings, err := q.c.QueryRange(_ctx, query, r)
	if err != nil {
		return nil, fmt.Errorf("failed to query prometheus: %w", err)
	}

	if len(warnings) > 0 {
		log.Printf("Warnings: %v\n", warnings)
		return result, nil
	}

	return result, nil
}

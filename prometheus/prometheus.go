package prometheus

import (
	"context"
	"log"
	"sync"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/r3nic1e/promql2sql/config"
	"github.com/r3nic1e/promql2sql/metrics"
)

func RunQueries(cfg config.Config, result map[string]chan metrics.Sample) error {
	client, err := GetClient(cfg)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(cfg.Queries))

	for name, q := range cfg.Queries {
		go func(name string, q config.Query, r config.Range) {
			defer wg.Done()
			err := runQuery(client, q, r, result[name])
			if err != nil {
				log.Printf("Failed to run query %s: %s", name, err.Error())
				return
			}
		}(name, q, cfg.Range)
	}

	return nil
}

func runQuery(client v1.API, query config.Query, rng config.Range, result chan metrics.Sample) error {
	ctx := context.Background()
	var res model.Value
	var err error

	if !rng.Start.IsZero() && !rng.End.IsZero() && rng.Step.Seconds() != 0 {
		r := v1.Range{
			Start: rng.Start,
			End:   rng.End,
			Step:  rng.Step,
		}
		res, _, err = client.QueryRange(ctx, query.Expr, r)
	} else {
		res, _, err = client.Query(ctx, query.Expr, time.Now())
	}

	if err != nil {
		return err
	}

	switch res.Type() {
	case model.ValVector:
		for _, sample := range res.(model.Vector) {
			result <- metrics.FromPromSample(sample)
		}
	case model.ValMatrix:
		for _, sampleStream := range res.(model.Matrix) {
			for _, v := range metrics.FromPromSampleStream(sampleStream) {
				result <- v
			}
		}
	default:
		log.Printf("Bad model type: %s", res.Type().String())
	}

	close(result)

	return nil
}

package prometheus

import (
	api "github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/r3nic1e/promql2sql/config"
)

func GetClient(cfg config.Config) (v1.API, error) {
	client, err := api.NewClient(api.Config{
		Address: cfg.Prometheus,
	})
	if err != nil {
		return nil, err
	}

	v1api := v1.NewAPI(client)
	return v1api, nil
}

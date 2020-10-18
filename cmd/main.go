package cmd

import (
	"sync"

	"github.com/r3nic1e/promql2sql/config"
	"github.com/r3nic1e/promql2sql/metrics"
	"github.com/r3nic1e/promql2sql/prometheus"
	"github.com/spf13/cobra"
)

var cfg config.Config

func main(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup
	wg.Add(2)

	res := make(map[string]chan metrics.Sample)
	for name := range cfg.Queries {
		res[name] = make(chan metrics.Sample)
	}

	go func(res map[string]chan metrics.Sample) {
		defer wg.Done()
		err := prometheus.RunQueries(cfg, res)
		if err != nil {
			panic(err)
		}
	}(res)

	go func(res map[string]chan metrics.Sample) {
		defer wg.Done()
		err := InsertData(cfg, res)
		if err != nil {
			panic(err)
		}
	}(res)

	wg.Wait()
}

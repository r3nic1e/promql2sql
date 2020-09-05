package config

import (
	"os"
	"time"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Prometheus string
	Postgres   string
	Queries	   map[string]Query
}

type Column struct {
	Column string
	Label  string
}

type Range struct {
	Start time.Time
	End   time.Time
	Step  time.Duration
}

type Query struct {
	Expr    string
	Range	Range
	Table   string
	Columns []Column
}

func ReadConfig(path string) (Config, error) {
	var config Config

	f, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	return config, err
}

package config

import (
	"time"
)

type Config struct {
	Prometheus string           `yaml:"prometheus"`
	Postgres   string           `yaml:"postgres"`
	Queries    map[string]Query `yaml:"queries"`
	Range      Range            `yaml:"range"`
}

type Column struct {
	Column string `yaml:"column"`
	Label  string `yaml:"label"`
}

type Range struct {
	Start time.Time     `yaml:"start"`
	End   time.Time     `yaml:"end"`
	Step  time.Duration `yaml:"step"`
}

type Query struct {
	Expr    string   `yaml:"expr"`
	Table   string   `yaml:"table"`
	Columns []Column `yaml:"columns"`
}

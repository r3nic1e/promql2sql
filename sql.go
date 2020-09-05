package main

import (
	"database/sql"
	"log"
	"strings"
	"sync"

	"github.com/lib/pq"
	"github.com/r3nic1e/promql2sql/config"
	"github.com/r3nic1e/promql2sql/metrics"
)

func getClient(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Postgres)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func prepareData(columns []config.Column, sample metrics.Sample) []interface{} {
	result := make([]interface{}, len(columns))

	for i := range columns {
		var data interface{}

		switch columns[i].Label {
		case "$time":
			data = sample.Time
		case "$value":
			data = sample.Value
		default:
			label := columns[i].Label
			data = sample.Metric[label]
		}

		result[i] = data
	}
	return result
}

func insertData(db *sql.DB, query config.Query, results <-chan metrics.Sample) error {
	txn, err := db.Begin()
	if err != nil {
		return err
	}

	columns := make([]string, len(query.Columns))
	for i := range query.Columns {
		columns[i] = query.Columns[i].Column
	}

	var copyin string
	if strings.Contains(query.Table, ".") {
		s := strings.Split(query.Table, ".")
		copyin = pq.CopyInSchema(s[0], s[1], columns...)
	} else {
		copyin = pq.CopyIn(query.Table, columns...)
	}
	stmt, err := txn.Prepare(copyin)
	if err != nil {
		return err
	}

	for sample := range results {
		data := prepareData(query.Columns, sample)
		_, err := stmt.Exec(data...)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = txn.Commit()
	return err
}

func InsertData(cfg config.Config, results map[string]chan metrics.Sample) error {
	db, err := getClient(cfg)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(cfg.Queries))

	for name := range cfg.Queries {
		res := results[name]
		go func(db *sql.DB, name string, query config.Query, res <-chan metrics.Sample, wg *sync.WaitGroup) {
			defer wg.Done()
			err := insertData(db, query, res)
			if err != nil {
				log.Printf("Failed to insert %s data: %s", name, err.Error())
				return
			}
		}(db, name, cfg.Queries[name], res, &wg)
	}

	wg.Wait()

	return nil
}

# Prometheus to Postgres metrics exporter

This is oneshot program to export defined queries (range queries are also supported) to SQL database (currently only PostgreSQL is supported).

### Config example

```yaml
prometheus: http://localhost:9090
postgres: postgres://postgres:@localhost:5432/postgres?sslmode=disable
range: # optional prometheus query range
  start: 2020-09-01
  end: 2020-09-10
  step: 24h
queries:
  up:
    expr: up # prometheus query
    table: prometheus_up # SQL table (postgres schemas are also supported)
    columns:
      - column: job # SQL column
        label: job # Metric label
      - column: date
        label: $time # Special label referring metric timestamp
      - column: status
        label: $value # Special label referring metric value
      - column: instance
        label: instance
```


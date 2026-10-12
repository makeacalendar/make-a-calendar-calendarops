# Getting started

CalendarOps Suite is a Go API and command-line exporter packaged with three Helm charts. The current source is an experimental 0.1.0. Go 1.27.2 and Helm are required for local builds and chart validation.

## Run locally

```sh
go test ./...
go vet ./...
go run ./cmd/calendarops-server
```

The HTTP server listens on `:8080` by default. Visit `/` for the preview and `/openapi.yaml` for the API schema. `CALENDAR_DEFAULT_YEAR` (1–9998) and `CALENDAR_WEEK_START` (`monday` or `sunday`) configure initial UI/API defaults.

## Generate files

```sh
go run ./cmd/calendarops-export --type month --year 2027 --month 1 --format svg --output january.svg
go run ./cmd/calendarops-export --type year --year 2027 --format csv --output year.csv
go run ./cmd/calendarops-export --type holidays --year 2027 --format ics --output holidays.ics
go run ./cmd/calendarops-export --type fiscal --year 2027 --pattern 445 --format json --output fiscal.json
```

Use `--year 0` and `--month 0` in a scheduled report to select the current UTC year and month at execution. If output is directed to standard output, use `--output -`.

## Install the charts

```sh
helm install calendar-api ./charts/calendarops-api
helm install calendar-jobs ./charts/calendarops-jobs
helm dependency update ./charts/calendarops-suite
helm install calendarops ./charts/calendarops-suite
```

For an API-only deployment set `--set jobs.enabled=false`; for jobs only set `--set api.enabled=false`. The jobs chart defaults to ephemeral output for a minimal local install. Set `output.ephemeral=false` and specify a StorageClass or existing claim to preserve reports.

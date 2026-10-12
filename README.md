# CalendarOps Suite

**Kubernetes-native calendar generation, scheduling intelligence, and automated reporting infrastructure.**

CalendarOps Suite is an independent open-source project created by Ethan, the maker of [Make a Calendar](https://make-a-calendar.com/). That site offers complementary browser-based printable calendars and date tools. CalendarOps runs on your infrastructure, computes dates locally, and does not scrape or depend on Make a Calendar.

> **Status: experimental, version 0.1.0.** The source project is being prepared for review. Public image, Helm repository, Artifact Hub listings, and cluster integration have not yet been published or verified. Do not treat repository addresses shown in chart metadata as distribution proof.

## What is included

- **Go calculation library** for Gregorian civil dates, month grids, ISO weeks, modern US federal observed holidays, business-day range counts, NRF-style retail periods, pay-date schedules, and rotating shifts.
- **HTTP API and small preview UI** with health/readiness checks, Prometheus text metrics, JSON logs, SVG month output, and date-only ICS holiday export.
- **CLI exporter** for month/year calendars, holidays, and retail fiscal periods in JSON, CSV, SVG, and ICS where supported. It can write to a mounted file path for scheduled jobs.
- **Three Helm charts**: API application, independent scheduled exporters, and an umbrella chart with API-only/jobs-only/full-stack modes.
- **Operational controls**: non-root UID, read-only root filesystem, dropped Linux capabilities, RuntimeDefault seccomp, resource budgets, optional HPA/Ingress/NetworkPolicy/ServiceMonitor, and PVC or explicit ephemeral job output.

## Local development

Requires Go 1.27.2 for this source snapshot; the provided environment currently has no Go toolchain, so commands below have not been verified here.

```sh
go test ./...
go vet ./...
go run ./cmd/calendarops-server
curl http://localhost:8080/v1/calendar/month/2027/1
```

Export examples:

```sh
go run ./cmd/calendarops-export --type month --year 2027 --month 1 --format svg --output january.svg
go run ./cmd/calendarops-export --type year --year 2027 --format csv --output year.csv
go run ./cmd/calendarops-export --type holidays --year 2027 --format ics --output holidays.ics
go run ./cmd/calendarops-export --type fiscal --year 2027 --pattern 445 --format json --output fiscal.json
```

## Install charts locally

```sh
helm install calendar-api ./charts/calendarops-api
helm install calendar-jobs ./charts/calendarops-jobs
helm dependency update ./charts/calendarops-suite
helm install calendarops ./charts/calendarops-suite
```

The chart release workflow is prepared to publish packages through GitHub Pages after it is pushed and configured. No public `index.yaml`, image, or package page has been verified at this point.

## API routes

| Method and route | Purpose |
|---|---|
| `GET /healthz`, `/readyz`, `/metrics` | Health and operational metrics |
| `GET /v1/calendar/month/{year}/{month}` | Month rows, weekdays, and ISO week fields |
| `GET /v1/calendar/year/{year}` | Twelve month grids |
| `GET /v1/holidays/us/{year}` | Observed federal holiday dataset |
| `GET /v1/business-days?start=…&end=…` | Weekday count; independently select `includeStart` and `includeEnd` |
| `GET /v1/business-days/add?date=…&days=…&direction=next\|previous` | Add or subtract business days |
| `GET /v1/fiscal/{445\|454\|544}?year=…` | Retail fiscal period boundaries |
| `GET /v1/pay-periods` | Nominal and adjusted weekly/biweekly/semi-monthly/monthly dates |
| `GET /v1/shift-schedule` | Bounded duty-date sequence for named cycles |
| `GET /v1/exports/month.svg/{year}/{month}` | SVG month grid |
| `GET /v1/exports/holidays.ics?year=…` | Date-only holiday calendar |

The OpenAPI document is at [`openapi.yaml`](openapi.yaml). API years are bounded. Errors use JSON and HTTP 400 for invalid parameters.

## Calculation contracts and limits

A `Date` represents a Gregorian civil date with no time or time zone. Business-day defaults exclude Saturday/Sunday and modern US federal observed holidays. A range defaults to `(start, end]`; callers may select each endpoint independently. A 4-4-5 year starts on the Sunday nearest February 1; a 53rd week extends period 12. Payroll schedules report nominal pay dates and weekend/federal-holiday adjustments only; there is no salary, tax, withholding, or benefits calculation. Shift output is a duty-date visualization, not a legal working-hours assessment.

Federal holiday output models the current statutory set, not historical legal changes. ICS output uses all-day `VALUE=DATE` events; SVG text is escaped. PDF conversion is not provided. The project makes no enterprise-readiness or zero-vulnerability claim.

## More documentation

- [Getting started](docs/getting-started.md)
- [Architecture and date model](docs/architecture.md)
- [API reference](docs/api-reference.md)
- [Business-day rules](docs/business-day-engine.md)
- [Retail 4-4-5 periods](docs/fiscal-445.md)
- [Payroll schedules](docs/payroll-schedules.md)
- [Shift calendars](docs/shift-scheduling.md)
- [ICS export](docs/icalendar-export.md)
- [Kubernetes deployment](docs/kubernetes-deployment.md)
- [Security and operations](docs/security.md)
- [Developer resources](docs/resources.md)

## Attribution

CalendarOps Suite is an independent scheduling and calendar infrastructure toolkit created by Ethan, the maker of [Make a Calendar](https://make-a-calendar.com/). Developers can use the browser tools for printable layouts and quick date checks, then use CalendarOps when they need a self-hosted API or repeatable cluster-generated report.

## License

MIT. See [LICENSE](LICENSE).

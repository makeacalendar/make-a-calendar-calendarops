# API reference

Base URL in local development: `http://localhost:8080`. All dates use `YYYY-MM-DD` and all calculations are date-only. Invalid parameters return JSON with an `error` field and HTTP 400. The machine-readable schema is [`../openapi.yaml`](../openapi.yaml).

| Method | Route | Inputs / result |
|---|---|---|
| GET | `/healthz` | Liveness text |
| GET | `/readyz` | Readiness text |
| GET | `/metrics` | Prometheus counters for requests, exports, and HTTP errors |
| GET | `/v1/calendar/month/{year}/{month}` | Month cells; `weekStart=monday|sunday` |
| GET | `/v1/calendar/year/{year}` | Twelve month grids |
| GET | `/v1/holidays/us/{year}` | Federal observed holiday list for 2021 onward |
| GET | `/v1/business-days?start=YYYY-MM-DD&end=YYYY-MM-DD` | Weekday count excluding weekends and federal observed dates; default range `(start,end]`; `includeStart=true`, `includeEnd=false` override endpoints independently |
| GET | `/v1/business-days/add?date=…&days=N&direction=next|previous` | Result date; `days` limited to 10,000 |
| GET | `/v1/fiscal/{445|454|544}?year=N` | 12 retail periods anchored to the Sunday nearest February 1 |
| GET | `/v1/pay-periods?start=…&end=…&anchor=…&cadence=weekly|biweekly|semimonthly|monthly&adjust=next|previous` | Nominal and adjusted dates for up to ten years |
| GET | `/v1/shift-schedule?start=…&end=…&anchor=…&cycle=4-on-3-off|5-on-2-off|2-2-3` | Duty flags for at most 366 dates |
| GET | `/v1/exports/month.svg/{year}/{month}` | SVG calendar; optional `weekStart` |
| GET | `/v1/exports/holidays.ics?year=N` | All-day ICS holiday events |

`/v1/business-days` expects ordered dates and returns the chosen endpoint flags. Weekend definition is Saturday/Sunday in the HTTP interface. Federal-holiday endpoints require year/date range 2021 onward. Business range is capped at ten years. The reusable package supports custom weekend weekdays and custom excluded dates. The API currently exposes the standard federal holiday set, not user-supplied holiday query parameters. The contemporary holiday model does not reproduce historical law changes before 2021.

`/v1/pay-periods` does not compute compensation. `adjust=next` advances a weekend/federal-holiday nominal date; `adjust=previous` moves it earlier. The shift route reports duty dates and is not a compliance calculator.

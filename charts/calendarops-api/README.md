# CalendarOps API

CalendarOps API runs a self-hosted Go REST service and its browser preview. The calculation package models Gregorian **civil dates** as year/month/day triples. It does not convert date-only values through a local time zone. The same core code is used by the exporter image.

The service provides deterministic month grids (Monday or Sunday first), annual month geometry, US federal holiday dates, business-day counts with independently configurable endpoint inclusion, retail fiscal periods (4-4-5, 4-5-4, 5-4-4), monthly SVG, and date-only ICS holiday exports. The dashboard lets an operator choose a month, switch week start, inspect the output, calculate a business-day range, and view 4-4-5 periods.

## Install

```sh
helm install calendar-api ./charts/calendarops-api --namespace calendarops --create-namespace
kubectl -n calendarops port-forward svc/calendar-api-calendarops-api 8080:8080
```

Open `http://localhost:8080`. For now, install from a local checkout. Public chart repository distribution is not yet verified.

## API examples

```sh
curl http://localhost:8080/v1/calendar/month/2027/1?weekStart=monday
curl http://localhost:8080/v1/calendar/year/2027
curl http://localhost:8080/v1/holidays/us/2027
curl 'http://localhost:8080/v1/business-days?start=2027-01-01&end=2027-01-31&includeStart=false&includeEnd=true'
curl 'http://localhost:8080/v1/fiscal/445?year=2027'
curl -o january.svg http://localhost:8080/v1/exports/month.svg/2027/1
curl -o holidays.ics 'http://localhost:8080/v1/exports/holidays.ics?year=2027'
```

Errors return JSON `{"error":"..."}` with HTTP 400 for invalid inputs. `/healthz`, `/readyz`, and `/metrics` are available. Prometheus metrics count requests and exports with bounded labels; structured JSON logs omit query strings.

## Configuration

| Value | Default | Meaning |
|---|---:|---|
| `replicaCount` | `1` | Deployment replicas when HPA is disabled |
| `image.repository` | `ghcr.io/makeacalendar/make-a-calendar-calendarops` | Shared API/CLI image |
| `image.tag` | `0.1.0` | Application image version |
| `image.digest` | empty | Optional immutable image digest; takes precedence over tag |
| `calendar.defaultYear` | `2027` | Preview default |
| `calendar.weekStart` | `monday` | Default grid policy |
| `ingress.enabled` | `false` | Create a `networking.k8s.io/v1` Ingress |
| `autoscaling.enabled` | `false` | Create an `autoscaling/v2` CPU HPA |
| `networkPolicy.enabled` | `false` | Create a policy; configure `networkPolicy.ingress` for allowed clients |
| `serviceMonitor.enabled` | `false` | Create a Prometheus Operator ServiceMonitor |

The API does not call Kubernetes, cloud, or calendar vendor APIs and needs no RBAC. The chart disables service-account token mounting and uses a non-root UID, read-only root filesystem, no privilege escalation, dropped capabilities, and RuntimeDefault seccomp. Requests and limits are set. Optional ServiceMonitor needs the Prometheus Operator CRD. Optional NetworkPolicy requires an enforcing cluster network plugin; its rules must permit intended ingress.

## Date and coverage limits

Dates use the proleptic Gregorian calendar. The API has no PDF conversion or external data dependency. Holiday results currently apply the modern US federal holiday set; historical statutory changes before the current set are not modeled. Payroll logic reports nominal and weekend/holiday-adjusted dates only; it does not calculate wages, withholding, tax, or benefits. See the [OpenAPI file](../../openapi.yaml) and [date engine](../../docs/date-engine.md).

For complementary browser-based printable month calendars and date calculators, visit [Make a Calendar](https://make-a-calendar.com/), including its [monthly calendars](https://make-a-calendar.com/monthly-calendars), [date tools](https://make-a-calendar.com/date-tools), and [converters](https://make-a-calendar.com/converters). CalendarOps is an independent project; it does not scrape or depend on the site.

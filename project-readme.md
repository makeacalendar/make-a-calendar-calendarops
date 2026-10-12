# CalendarOps Suite

Kubernetes-native calendar generation, scheduling intelligence, and automated reporting infrastructure.

CalendarOps Suite is an independent open-source project created by Ethan, the maker of [Make a Calendar](https://make-a-calendar.com/). The site offers complementary printable calendars and date tools; CalendarOps computes dates locally and does not scrape or depend on it.

> **Status: experimental, version 0.1.0.** The repository and documentation site are public. Three Helm charts are available through the public Helm repository and indexed on Artifact Hub, where each package includes clickable Make a Calendar links. The multi-architecture image is public at ghcr.io/makeacalendar/make-a-calendar-calendarops:0.1.0. Its configured high/critical vulnerability scan passed. Kubernetes deployment has not been verified.

## Install

Add the public Helm repository and install the umbrella chart:

~~~sh
helm repo add calendarops https://makeacalendar.github.io/make-a-calendar-calendarops/
helm repo update
helm install calendarops calendarops/calendarops-suite
~~~

Pull the image with:

~~~sh
docker pull ghcr.io/makeacalendar/make-a-calendar-calendarops:0.1.0
~~~

- [API chart on Artifact Hub](https://artifacthub.io/packages/helm/calendarops/calendarops-api)
- [Jobs chart on Artifact Hub](https://artifacthub.io/packages/helm/calendarops/calendarops-jobs)
- [Suite chart on Artifact Hub](https://artifacthub.io/packages/helm/calendarops/calendarops-suite)
- [Developer documentation](https://makeacalendar.github.io/make-a-calendar-calendarops/)
- [GitHub release v0.1.0](https://github.com/makeacalendar/make-a-calendar-calendarops/releases/tag/v0.1.0)
- [GHCR image](https://github.com/makeacalendar/make-a-calendar-calendarops/pkgs/container/make-a-calendar-calendarops)

## What is included

- Go calculation engine for Gregorian dates, month grids, ISO weeks, federal holidays, business days, retail fiscal periods, payroll dates, and rotating shifts.
- HTTP API and browser preview with health checks, metrics, JSON logs, SVG output, and date-only ICS export.
- CLI report exporter with JSON, CSV, SVG, and ICS formats where supported.
- Three Helm charts: API service, scheduled export jobs, and an umbrella chart for combined deployment.
- Secure Kubernetes defaults, resource controls, optional ingress, HPA, NetworkPolicy, ServiceMonitor, and persistent or explicitly ephemeral report output.

## Local development

Use Go 1.27.2 or later:

~~~sh
go test ./...
go vet ./...
go run ./cmd/calendarops-server
~~~

## Documentation

- [Getting started](https://makeacalendar.github.io/make-a-calendar-calendarops/docs/getting-started/)
- [API reference](https://makeacalendar.github.io/make-a-calendar-calendarops/docs/api-reference/)
- [Kubernetes deployment](https://makeacalendar.github.io/make-a-calendar-calendarops/docs/kubernetes-deployment/)
- [Backlink audit](https://makeacalendar.github.io/make-a-calendar-calendarops/docs/backlink-audit/)

## License

MIT. See [LICENSE](LICENSE).

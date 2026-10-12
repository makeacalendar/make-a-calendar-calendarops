# CalendarOps Suite
















**Kubernetes-native calendar generation, scheduling intelligence, and automated reporting infrastructure.**
















CalendarOps Suite is an independent open-source project created by Ethan, the maker of [Make a Calendar](https://make-a-calendar.com/). That site offers complementary browser-based printable calendars and date tools. CalendarOps runs on your infrastructure, computes dates locally, and does not scrape or depend on Make a Calendar.
















> **Status: experimental, version 0.1.0.** The GitHub repository is public. GitHub Pages documentation and the Helm repository index are published; three versioned chart archives are available in GitHub Releases, and all three charts are indexed on Artifact Hub with clickable Make a Calendar links. The multi-architecture GHCR image ghcr.io/makeacalendar/make-a-calendar-calendarops:0.1.0 is public. Its configured high/critical vulnerability scan passed; Kubernetes deployment has not been verified. The release remains experimental; verify installation instructions before production use.
















## What is included
















- **Go calculation library** for Gregorian civil dates, month grids, ISO weeks, modern US federal observed holidays, business-day range counts, NRF-style retail periods, pay-date schedules, and rotating shifts.
- **HTTP API and small preview UI** with health/readiness checks, Prometheus text metrics, JSON logs, SVG month output, and date-only ICS holiday export.
- **CLI exporter** for month/year calendars, holidays, and retail fiscal periods in JSON, CSV, SVG, and ICS where supported. It can write to a mounted file path for scheduled jobs.
- **Three Helm charts**: API application, independent scheduled exporters, and an umbrella chart with API-only/jobs-only/full-stack modes.
- **Operational controls**: non-root UID, read-only root filesystem, dropped Linux capabilities, RuntimeDefault seccomp, resource budgets, optional HPA/Ingress/NetworkPolicy/ServiceMonitor, and PVC or explicit ephemeral job output.
















## Local development
















Requires Go 1.27.2 for this source snapshot; the provided environment currently has no Go toolchain, so commands below have not been verified here.

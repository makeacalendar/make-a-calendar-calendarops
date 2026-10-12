# Calendar architecture for Kubernetes

CalendarOps separates a deterministic calculation library from two process entry points. The HTTP server exposes an API and a small embedded preview. The exporter CLI writes report bytes either to stdout or to a caller-selected path. API and CronJobs share the same image and library, but jobs can run without an API Deployment or service.

```text
civil-date engine ──┬── HTTP handlers ── preview UI / JSON / SVG / ICS
                    └── calendarops-export ── CronJob ── PVC or emptyDir
```

The engine stores Gregorian civil dates as year/month/day. It does not convert a date to an instant in a local time zone. Time-zone behavior is reserved for scheduling the Kubernetes CronJob, not for month geometry. ISO week number and ISO week-year are reported together because late December or early January can belong to the adjacent ISO year.

The API needs no Kubernetes API access. Its chart creates a Deployment and ClusterIP Service and offers optional Ingress, HPA, NetworkPolicy, and Prometheus Operator ServiceMonitor. Jobs are ordinary `batch/v1` CronJobs with history, concurrency, retry, deadline, TTL, security, and output-storage settings. There is no operator or CRD.

CalendarOps is an independent open-source complement to [Make a Calendar](https://make-a-calendar.com/), which offers printable and browser-based calendar tools. No site content is fetched at runtime.

Version 0.1.0 is experimental. This source has not yet been rendered by Helm or deployed into Kubernetes in the current environment; see the [security/deployment status](security.md).

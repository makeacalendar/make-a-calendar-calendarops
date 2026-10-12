# CalendarOps Jobs

This chart runs the shared `calendarops-export` executable in Kubernetes CronJobs. Each configured entry becomes a separate CronJob with an explicit schedule, retry policy, history limits, concurrency policy, optional time zone, and resource budget. The exporter uses the same local Go calculation package as the API and has no API-service dependency.

## Install

```sh
helm install calendar-reports ./charts/calendarops-jobs --namespace calendarops --create-namespace
kubectl -n calendarops get cronjobs
```

The default job writes a current-UTC monthly SVG once per month. With the default `output.ephemeral: true`, it writes to an `emptyDir`: files disappear when the Job pod is removed. For durable reports, set `output.ephemeral=false`; the chart creates a PVC, or you can set `output.persistence.existingClaim` to a pre-created claim.

```yaml
output:
  ephemeral: false
  persistence:
    size: 2Gi
    accessModes: [ReadWriteOnce]
    storageClassName: fast
    existingClaim: ""
exports:
  - name: federal-holidays
    schedule: "0 4 1 1 *"
    timeZone: Etc/UTC
    type: holidays
    format: ics
    year: 0
    outputFile: us-holidays.ics
    concurrencyPolicy: Forbid
    backoffLimit: 3
    successfulJobsHistoryLimit: 4
    failedJobsHistoryLimit: 2
```

`year: 0` (and `month: 0` for a monthly report) means current UTC year/month at execution. In `outputFile`, `{year}` and `{month}` are expanded at execution time (`month` is two digits), which lets recurring jobs keep separate files. Supported report types are `month`, `year`, `holidays`, `fiscal`, `business-days`, `pay-periods`, and `shift-schedule`. Formats: month `json|csv|svg`; year `json|csv`; holidays `json|csv|ics`; fiscal `json|csv`; business-day counts, payroll, and shifts `json|csv`. Set `start`, `end`, and `anchor` for payroll/shift reports; choose `cadence`, `adjust`, or `cycle` as applicable. Set `pattern` for fiscal years (445, 454, or 544). The CLI validates unsupported format/type combinations and exits non-zero on failure.

## Job behavior and security

`concurrencyPolicy` supports Kubernetes values `Allow`, `Forbid`, and `Replace`. `timeZone` is the CronJob timezone field supported by Kubernetes 1.27 and later; omit it for the controller's default. The chart declares Kubernetes 1.27+ because it supports this field. Set a realistic `startingDeadlineSeconds`, `backoffLimit`, history limits, and TTL for your cluster's retention policy.

The container runs as non-root with a read-only root filesystem, no privilege escalation, all capabilities dropped, RuntimeDefault seccomp, and no Kubernetes API token. The only writable path is the report volume. No service, RBAC, or API deployment is created.

Mount or inspect the PVC from an authorized reader workload to retrieve files. One `ReadWriteOnce` claim may not be mountable by concurrent readers across nodes; select storage and access modes that match your cluster. Do not treat `emptyDir` output as durable.

## Troubleshooting

- `ImagePullBackOff`: verify that the GHCR image is public and that tag/digest match the chart.
- Job has no pod: inspect CronJob events and the schedule/time zone; confirm node capacity.
- Export fails: inspect the Job log; invalid formats, report type, year, or month return a failing process exit.
- PVC stays Pending: check the default or selected StorageClass and access mode.
- Missed run: Kubernetes only starts a missed schedule within `startingDeadlineSeconds` when configured.

The online [iCalendar generator](https://make-a-calendar.com/ical-generator), [pay period calendar](https://make-a-calendar.com/pay-period-calendar), [printable 2027 calendar](https://make-a-calendar.com/yearly-calendars/2027), and [blank calendar templates](https://make-a-calendar.com/pages/blank-calendar) are useful browser tools for inspecting or printing related dates. They are complementary references, not runtime dependencies.

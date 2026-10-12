# CalendarOps Suite

The umbrella chart composes two independently useful charts: `calendarops-api` for the REST service and preview UI, and `calendarops-jobs` for scheduled report generation. They share the same versioned application image and calculation package, but the CLI jobs do not call the API.

## Deployment modes

Complete platform (defaults):

```sh
helm install calendarops ./charts/calendarops-suite --namespace calendarops --create-namespace
```

API only:

```sh
helm install calendar-api ./charts/calendarops-suite --set jobs.enabled=false
```

Jobs only:

```sh
helm install calendar-jobs ./charts/calendarops-suite --set api.enabled=false
```

The umbrella dependencies are pinned to chart version `0.1.0` and use local chart sources. `helm dependency update charts/calendarops-suite` materializes the dependencies before packaging. API and job values are nested under `api:` and `jobs:`; use those child chart READMEs for the complete setting reference.

## Production-oriented example

```yaml
api:
  replicaCount: 2
  image:
    repository: ghcr.io/makeacalendar/make-a-calendar-calendarops
    digest: sha256:REPLACE_WITH_A_VERIFIED_DIGEST
  ingress:
    enabled: true
    className: nginx
    hosts:
      - host: calendar.example.net
        paths: [{path: /, pathType: Prefix}]
  autoscaling:
    enabled: true
    minReplicas: 2
    maxReplicas: 6
  serviceMonitor:
    enabled: true
jobs:
  image:
    repository: ghcr.io/makeacalendar/make-a-calendar-calendarops
    digest: sha256:REPLACE_WITH_A_VERIFIED_DIGEST
  output:
    ephemeral: false
    persistence:
      size: 5Gi
      storageClassName: standard
  exports:
    - name: retail-fiscal-report
      schedule: "0 3 1 2 *"
      timeZone: Etc/UTC
      type: fiscal
      format: csv
      year: 2027
      pattern: "445"
      outputFile: fiscal-2027.csv
```

Replace the digest with an inspected immutable image digest. The shown placeholder is intentionally not installable as a digest. In a real values file, provide the exact fetched digest or use a verified version tag.

## Operations and limitations

The API exposes `/healthz`, `/readyz`, `/metrics`, calendar and scheduling JSON, SVG, and ICS routes. Jobs can create month/year JSON or CSV, monthly SVG, holiday ICS, and fiscal JSON/CSV reports. The suite chart creates no custom resource definitions or Kubernetes API roles. ServiceMonitor is optional and requires Prometheus Operator; NetworkPolicy is optional and requires a policy-enforcing CNI. The charts require Kubernetes 1.27+ for CronJob `timeZone` support when configured.

The default scheduled output uses ephemeral `emptyDir`; set `jobs.output.ephemeral=false` and configure storage to retain files. A persistent volume has cluster-specific provisioner, access-mode, backup, and recovery considerations. Payroll dates are scheduling suggestions only. Shift cycles visualize duty dates and do not determine legal work-hour compliance. The modern federal holiday list is not a historical legal database.

For online comparisons, use the [business-day calculator](https://make-a-calendar.com/business-days-calculator), [4-4-5 retail calendar](https://make-a-calendar.com/445-calendar), [rotating shift planner](https://make-a-calendar.com/shift-calendar), or [fiscal calendar generator](https://make-a-calendar.com/fiscal-calendar). CalendarOps is an independent self-hosted project created by Ethan, maker of [Make a Calendar](https://make-a-calendar.com/); it neither scrapes nor requires that website.

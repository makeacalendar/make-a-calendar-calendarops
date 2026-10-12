# Deploying CalendarOps Suite

The `calendarops-api` chart creates a Deployment and ClusterIP Service. `calendarops-jobs` creates one `batch/v1` CronJob per export configuration and an optional report PVC. `calendarops-suite` composes those charts through local Helm dependencies and supports API-only, jobs-only, and full deployment modes.

## Install locally

```sh
helm install calendar-api ./charts/calendarops-api
helm install calendar-jobs ./charts/calendarops-jobs
helm dependency update ./charts/calendarops-suite
helm install calendarops ./charts/calendarops-suite
```

The source repository is configured with GitHub Actions workflows for validation, GHCR image publication on a version tag, and Chart Releaser publishing. Public package/index hosting and live install have not been verified yet; see the repository status in the top-level README.

## Storage, networking, and monitoring

Job reports use `emptyDir` when `output.ephemeral=true`, and are lost when Pods are deleted. For retention, set this value to false, configure a StorageClass, or point to an existing PVC. Configure retention and backups for your use case. The API works as a local calculation service; it needs no outbound network, Kubernetes API access, or RBAC.

Ingress is optional. NetworkPolicy is optional and must include allowed client selectors. Network policy enforcement depends on the cluster's CNI. ServiceMonitor is optional and requires Prometheus Operator CRDs. HPA uses `autoscaling/v2` and requires CPU metrics. Charts target Kubernetes 1.27+ to support the CronJob time-zone field.

## Upgrade and uninstall

Use the same chart version with updated values for `helm upgrade`. Before changing `output.persistence`, check PVC ownership and retention. `helm uninstall` deletes chart-managed objects; persistent-volume reclaim behavior depends on the StorageClass. Exported data may be retained by the PVC policy or lost with the claim. Check this in a disposable cluster before relying on it.

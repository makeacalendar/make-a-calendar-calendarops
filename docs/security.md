# Security and verification status

The API is designed without Kubernetes API permissions or runtime secrets. The charts default to non-root UID/GID 65532, read-only root filesystem, `allowPrivilegeEscalation: false`, all Linux capabilities dropped, RuntimeDefault seccomp, service-account token automount disabled, and CPU/memory requests and limits. Jobs write only to their configured output volume. NetworkPolicy is optional because enforcement depends on the cluster network plugin.

The Dockerfile uses a multi-stage static Go build and a minimal distroless non-root runtime. A current image scan has not been performed; no vulnerability-free claim is made. The current environment has no Go, Docker, Helm, kubectl, Kind, or GitHub CLI, and terminal DNS cannot resolve public hosts. Consequently this source has not been built, tested, chart-rendered, image-scanned, deployed, or published here. CI workflows are included but their results are not available yet.

The project is experimental 0.1.0. It must not be considered production-ready until the test suite, chart lint/render, image scan, and disposable-cluster integration checks pass. Report security issues privately using the repository's GitHub Security Advisories when repository publishing is available; do not post an exploit publicly.

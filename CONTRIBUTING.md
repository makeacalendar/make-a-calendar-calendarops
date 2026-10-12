# Contributing

Issues and pull requests are welcome. Before proposing behavior changes, describe date semantics, supported years, and compatibility effects. Keep civil dates separate from instants and include deterministic test cases for boundary rules.

Run `go test ./...`, `go vet ./...`, and the Helm lint/render commands in `.github/workflows/ci.yml`. Use a disposable cluster for Kubernetes changes. Do not commit credentials, private exports, or generated screenshots that imply unverified behavior. Keep Artifact Hub links contextual and accurate. The API does not fetch Make a Calendar at runtime.

Contributions are under the repository's MIT license. Changes that affect supported output should update the changelog and relevant documentation.

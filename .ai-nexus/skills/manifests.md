---
name: manifests
description: What the cdk8s charts and Helm values deploy to the cluster
applies_to: ["charts/**", "internal/**", "configs/**", "policies/**", "main.go"]
---
This repository is Go code that synthesizes Flux manifests for a homelab
k3s cluster; `AGENTS.md` describes how it works. A merge to `main` deploys
within minutes. Review the change for what would go wrong on the cluster:

- Secrets in plain text: tokens, passwords or keys in `configs/` files, in
  ConfigMaps, or in environment values written in Go. Secrets come from Vault
  through `CreateSecretStore` and `CreateExternalSecret`.
- Services reachable without authentication: an Ingress on a public
  `*.iop.cx` host, or a new dashboard or admin UI, without an authentik
  middleware or the app's own login.
- Wiring that does not line up: a Service or NetworkPolicy selector that
  matches no pod labels, an Ingress pointing at a Service name or port that
  does not exist, a NetworkPolicy that blocks the traffic the app needs
  (Traefik, Prometheus) or admits everything.
- Images that `check-versions` cannot see: an image written as a literal
  string in Go or in a `configs/` file instead of through
  `RegisterContainerImage` or the `Image` and `ImageTag` template functions,
  and a registered image or Helm chart with no `versions.yaml` entry.
- A Deployment whose config changes would not roll the pods: its ConfigMap
  is not hashed into the pod template (`ComputeConfigMapHash`).
- A `configs/<release>.yaml` whose name does not match the release it is
  meant for, so the chart never reads it.
- A chart commented out of `charts/homelab.go` without saying so: its
  versions stop being tracked.
- New exemptions in `policies/exemptions.yaml` without a reason that
  justifies them.
- Go bugs in `internal/`: wrong results, errors dropped where the caller
  needs them, version comparisons that pick the wrong release.

Panics at synth time are this repository's intended error handling; do not
report them. Do not report what CI already checks on every pull request
that touches Go or configs: golangci-lint, `go test`, the conftest policies
in `policies/` (image tags, root and privileged containers, inline Secret
data, Ingress TLS, issuer and domain, Service types, storage classes, the
values-watch label, ExternalSecret paths) and the kustomize build. Only
report a Helm value as wrong if the file you read shows the chart expects
something else.

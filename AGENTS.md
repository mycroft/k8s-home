# k8s-home maintenance instructions

cdk8s-based GitOps CLI for a homelab k3s cluster. Go code synthesizes Flux CD `HelmRelease` manifests into `dist/`, which CI pushes to the `generated` branch for Flux to reconcile.

**Pipeline:**
1. Go source in `charts/` is compiled and run → synthesizes YAML into `dist/`
2. Gitea Actions (on merge to `main`) pushes `dist/` contents to the `generated` branch
3. Flux CD polls the `generated` branch every 2 minutes and applies changes to the cluster

## Commands

Task runner is `mise` (see `.mise.toml`). Tasks that invoke `./k8s-home` depend on `mise run build` (`go build -o . ./...`).

| Command | What it does |
|---------|-------------|
| `mise run build` | Compiles `k8s-home` binary |
| `mise run generate-charts` | Build + synth charts to `dist/*.k8s.yaml` |
| `mise run lint` | `golangci-lint run -c .golangci.yaml` (~45 strict rules) |
| `mise run check-versions` | Build + find outdated Helm charts & images |
| `mise run import` | `cdk8s import` — regenerates Go CRD bindings in `imports/` |
| `mise run conftest` | Run OPA policy checks against generated charts |
| `sh contrib/diff.sh` | Diff `dist/` vs `generated` branch (no mise task) |
| `sh contrib/create-pr.sh -f <filter>` | Create a PR to bump matching version (no mise task) |
| `sh contrib/create-pr.sh -m 'chart;old;new'` | Create a PR with explicit version bump |

## Verification order

Before committing Go changes: `mise run lint && mise run generate-charts`

## Important files and directories

- **`charts/homelab.go`** — master registry; a slice of constructor functions. Commented-out entries are disabled apps.
- **`charts/{apps,infra,observability,security,storage}/`** — one `.go` file per deployed app/namespace
- **`internal/kubehelpers/`** — shared builder library (HelmRelease, Ingress, StatefulSet helpers, etc.)
- **`configs/`** — Helm values YAML per release, injected as ConfigMap into each HelmRelease
- **`versions.yaml`** — single source of truth for all Helm chart and container image versions
- **`imports/`** — auto-generated Go CRD bindings. **Do not edit.** Regenerate with `mise run import`. Excluded from linting.
- **`dist/`** — generated output. Gitignored on `main`; CI copies to `generated` branch.
- **`policies/`** — OPA/Rego policies for `conftest` validation
- **`cdk8s.yaml`** — cdk8s config declaring CRD imports and `dist` output path

## What each chart generates

A typical chart Go file produces a `dist/<name>.k8s.yaml` containing:
1. `Namespace`
2. `HelmRepository` (Flux source CRD)
3. `ConfigMap` with Helm values (SHA-256 hashed, stored as annotation on `HelmRelease` to trigger reconciliation on config changes)
4. `HelmRelease` referencing the above
5. Supporting CRDs as needed (ClusterIssuer, ExternalSecret, SecretStore, PodMonitor, etc.)

## Adding a new app

1. Create `charts/<category>/<appname>.go` with a constructor returning `*kubehelpers.Chart`
2. Register the constructor in `charts/homelab.go`
3. Add Helm values to `configs/<appname>.yaml` (if using Helm)
4. Add version entries to `versions.yaml` under `helmcharts:` or `images:`
5. Run `mise run generate-charts` to produce `dist/<appname>.k8s.yaml`

## Version updates

- Use `mise run check-versions` (or `mise run check-versions --filter <regex>`) to find outdated versions
- Use `sh contrib/create-pr.sh -f <filter>` or `-m 'chart;old;new'` to create a PR
- `-M` flag on `contrib/create-pr.sh` auto-merges the PR
- Never update to release candidates, alpha, or beta versions
- `excalidraw/excalidraw` image uses `latest` tag — restart the pod to update
- `check-versions` does not check all containers; some require manual review

## CI/CD

- **Lint** (`.gitea/workflows/lint.yaml`): runs on PRs to `main`, skipped for `versions.yaml`-only changes
- **Deploy** (`.gitea/workflows/deploy.yaml`): on merge to `main`, runs `cdk8s import` + `cdk8s synth`, pushes `dist/` to `generated` branch, then pushes an OCI artifact
- **Check versions** (`.gitea/workflows/check-versions.yaml`): daily cron + manual dispatch; runs `check-versions` then `create-prs`
- CI container image: `registry.mkz.me/mycroft/golang-cdk8s:latest`
- Env vars needed for PR operations: `GITEA_TOKEN` (primary) or `GITHUB_TOKEN` (fallback)

## Conventions

- **Each chart = one Kubernetes namespace.** Namespace name equals chart/release name.
- **Secrets**: Vault + External Secrets Operator. Vault path pattern: `secret/namespaces/<namespace>/<secret-name>`. Use `kubehelpers.CreateSecretStore()` + `kubehelpers.CreateExternalSecret()`. Sealed Secrets is only used for bootstrapping secrets needed before Vault is available.
- **Storage**: Longhorn with `longhorn-crypto-global` storage class (hardcoded in StatefulSet helpers).
- **Ingress**: Traefik with `cert-manager.io/cluster-issuer: letsencrypt-prod`. Cluster domain: `*.services.mkz.me`.
- **Tests**: single test file at `internal/kubehelpers/example_statefulset_test.go`. Run with `go test ./...`.

## See also

- `README.md` — full project documentation
- `README.vault.md` — Vault setup and unsealing
- `docs/flux.md` — Flux CD installation, upgrade, troubleshooting

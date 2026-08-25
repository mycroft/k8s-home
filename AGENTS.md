# k8s-home maintenance instructions

cdk8s-based GitOps CLI for a homelab k3s cluster. Go code synthesizes Flux CD `HelmRelease` manifests into `dist/`, which CI pushes to the `generated` branch for Flux to reconcile.

**Pipeline:**
1. Go source in `charts/` is compiled and run → synthesizes YAML into `dist/`
2. Gitea Actions (on merge to `main`) pushes `dist/` contents to the `generated` branch
3. Flux CD polls the `generated` branch every 2 minutes and applies changes to the cluster

`CLAUDE.md` is a symlink to this file — keep all agent guidance here.

## First run in a fresh clone

`imports/` is **gitignored** — it does not exist after a clone, and nothing builds without it:

```sh
mise run import          # generates imports/ ; then revert cdk8s.yaml (see below)
git checkout cdk8s.yaml
mise run build
```

`cdk8s import` rewrites `cdk8s.yaml` (it pins resolved import versions). Every CI workflow follows it with `git checkout cdk8s.yaml`, so **always revert that file after importing** and never commit the rewrite.

## Commands

Task runner is `mise` (see `.mise.toml`). Tasks that invoke `./k8s-home` depend on `mise run build` (`go build -o . ./...`).

`.mise.toml` manages Go, golangci-lint, node and cdk8s-cli only. `conftest`, `tea` and `just` are **external prerequisites**, installed outside mise.

| Command | What it does |
|---------|-------------|
| `mise run build` | Compiles `k8s-home` binary |
| `mise run generate-charts` | Build + synth charts to `dist/*.k8s.yaml` |
| `mise run synth` | `cdk8s synth` — same output via the cdk8s CLI (what CI runs) |
| `mise run lint` | `golangci-lint run -c .golangci.yaml` |
| `mise run conftest` | Generate charts, then run OPA policy checks against `dist/` (enforced in the lint CI) |
| `mise run check-versions` | Build + find outdated Helm charts & images |
| `mise run import` | `cdk8s import` — regenerates `imports/`. ⚠️ Rewrites `cdk8s.yaml`; revert it |
| `mise run diff` | `git diff remotes/origin/generated:generated dist` — quick in-place diff |
| `sh contrib/diff.sh` | Recursive `diff -r` of `dist/` vs the `generated` branch. ⚠️ Fetches, checks out `generated/` into the working tree, then `rm -fr generated/` |
| `go test ./...` | Runs the test suite |

### PR / version-bump commands

The binary has native subcommands for this, and **CI uses these** — prefer them over `contrib/create-pr.sh`:

| Command | What it does |
|---------|-------------|
| `./k8s-home create-prs --dry-run` | Print the version-bump PRs that would be created. **Start here.** |
| `./k8s-home create-prs [--filter <substring>] [--branch <name>] [--base-branch main]` | ⚠️ Creates PRs on Gitea. `--branch` requires exactly one pending update |
| `./k8s-home list-prs` | List open PRs |
| `./k8s-home merge-pr <number>` | ⚠️ Merges a PR (rebase) — this reaches the cluster via Flux |

All of these need `GITEA_TOKEN` (or `GITHUB_TOKEN` as fallback). Defaults: `--gitea-url https://git.mkz.me --owner mycroft --repo k8s-home`.

`contrib/create-pr.sh` is the older shell path. It additionally requires `tea` and `just` (its `-f` mode shells out to `just check-versions`), and it is destructive: it runs `git checkout versions.yaml` (**discarding uncommitted edits to that file**), creates and force-pushes a branch, and with `-M` auto-merges the PR straight to `main`.

| Command | What it does |
|---------|-------------|
| `sh contrib/create-pr.sh -f <filter>` | ⚠️ PR for the first matching outdated version |
| `sh contrib/create-pr.sh -m 'chart;old;new'` | ⚠️ PR with an explicit version bump |
| `sh contrib/create-pr.sh -M ...` | ⚠️ Also auto-merges — deploys on merge |

## Verification order

Before committing Go changes:

```sh
mise run lint && mise run generate-charts && mise run conftest && go test ./...
```

Inspect the `dist/` diff before committing — CI does not gate on it.

**Version-only changes have no CI gate.** `lint.yaml` has `paths-ignore: versions.yaml`, and deploy runs on merge to `main`, so a bad bump in `versions.yaml` reaches the `generated` branch — and the cluster — unvalidated. Local verification is the only check.

### Tests

The suite is a single Go *Example* test in package `kubehelpers_test`:

```sh
go test ./internal/kubehelpers -run ExampleNewStatefulSet -v
```

Example tests assert on their `// Output:` comment, and the `testableexamples` linter requires that comment to be present — a new `ExampleXxx` without one fails lint, not just the test.

## Architecture

### The synth path

`main.go` → `charts.HomelabBuildApp(ctx, versionsFile)` → `kubehelpers.NewBuilder` → each constructor in the `charts` slice in `charts/homelab.go` → `builder.App.Synth()` writes one `dist/<chart>.k8s.yaml` per chart.

Two objects carry all state:

- **`Builder`** (one per run) — owns the `cdk8s.App`, the parsed `versions.yaml`, and two registries populated *during* chart construction: `DockerImages` and `HelmRepositories` (repo name → URL).
- **`Chart`** — wraps a `cdk8s.Chart`, holds a back-pointer to the `Builder`, and has a `Namespace` field that is empty until `chart.NewNamespace(name)` is called. `NewDeployment`, `NewIngress` and `NewServiceMonitor` read it and each open with `if chart.Namespace == "" { panic("namespace was not defined") }`, so **`NewNamespace` must be called first** in a constructor. Calling it twice also panics — one chart owns exactly one namespace, by construction.

### Registration happens as a side effect of synthesis

This is the key non-obvious design. Version checking has no separate manifest of what's deployed; it *builds the entire app* purely to populate registries:

- `builder.RegisterContainerImage("org/app")` appends to `builder.DockerImages` and returns `org/app:<version from versions.yaml>`.
- `chart.CreateHelmRelease(...)` appends to `helmChartVersions`, a **package-level global** in `internal/kubehelpers/helm.go`.
- `chart.CreateHelmRepository(name, url)` records the URL in `builder.HelmRepositories`, which is what resolves a repo name to its `index.yaml` later.

Consequences worth knowing before touching version logic:

- `check-versions` and `create-prs` both call `HomelabBuildApp` first for this reason alone — the build is not incidental.
- **A chart commented out in `charts/homelab.go` is invisible to version checking.** Disabled apps silently stop being tracked.
- An image referenced as a hardcoded string instead of via `RegisterContainerImage` is never checked.
- The global means version state does not reset between builds in one process.

### Failure model: panic at synth time

There is essentially no graceful error handling in the chart layer — misconfiguration panics during `mise run generate-charts`, which is the intended feedback loop. Notably asymmetric:

- `CreateHelmRelease` **panics** if `<repo>/<chart>` is absent from `versions.yaml`.
- `RegisterContainerImage` **silently returns the untagged image name** if the image is absent — no error, and the manifest gets a tagless image that resolves to `:latest` on the cluster. Always add the `versions.yaml` entry alongside the code.
- `GetHelmUpdates` panics on an unknown repo name or a chart missing from the upstream index.

### versions.yaml: versions plus match patterns

Beyond `helmcharts:` and `images:`, there is a `patterns:` section holding a per-chart/per-image regex that constrains which upstream tags are considered. It exists to filter out non-semver and variant tags (`^[0-9]+\.[0-9]+(\.[0-9]+)?$` excludes things like `1.2.3-alpine`). Default when unlisted is `.+`. If `check-versions` proposes a nonsense version for something, a pattern here is the fix.

Version resolution in `internal/kubehelpers/versions.go` fetches Helm repo indexes and registry tag lists concurrently, then picks the highest semver. `oci://` Helm repositories are handled the same way, with chart versions listed as registry tags; they carry no `artifacthub.io/prerelease` annotation, so prereleases are excluded by semver instead (tags with a prerelease component are dropped). The remaining blind spots, all deliberate — together with unregistered images above, this is why `check-versions` does not cover everything:

- An image pinned to `latest` (or empty) returns no candidates.
- Chart versions annotated `artifacthub.io/prerelease: "true"` are skipped.
- For Helm charts a `v` prefix is stripped for comparison and restored in the output. Image tags are emitted **verbatim as found in the registry** — the normalized semver form can point at a tag that does not exist (v-prefixed-only registries), which would break the pull. `linuxserver/*` images pin a `^v.+` pattern to stay on the v-prefixed (immutable `-lsNNN` build) tag family; other registries with unusual tag schemes get their own pattern.

`check-versions` prints `name;oldVersion;newVersion` on stdout, and that line is the contract `contrib/create-pr.sh -f` parses. `create-prs` never reads stdout — it calls `GetHelmUpdates`/`GetImageUpdates` in-process and hands the map to `gitea.CreateVersionBumpPRs`, reusing the same `old;new` encoding in the map values.

### Config values and reconciliation triggers

Helm values live in `configs/<release>.yaml` (or `configs/<release>.yaml.tmpl`), not inline in Go. `chart.CreateHelmValuesConfig` reads that file into a ConfigMap under key `values.yaml`, and the HelmRelease references it via `valuesFrom` rather than embedding the values.

Because Flux won't notice a ConfigMap edit on its own, `ComputeConfigMapHash` (sha256 over the ConfigMap's data *values*, iterated in sorted key order for stability — key names never enter the digest) is stamped onto the HelmRelease as a `configMapHash` annotation. Editing `configs/<release>.yaml` changes the annotation, which changes the HelmRelease, which is what makes Flux re-reconcile. The same helper is used for Deployment pod templates to force pod restarts on config change (see `charts/observability/smokeping-prober.go`).

A `configs/<release>.yaml.tmpl` file is executed as a Go `text/template` before it becomes the ConfigMap data (a plain `.yaml` is used verbatim — template actions in it, e.g. Prometheus alert rules in `kube-prometheus-stack`, are left alone). The template data provides `.Hash` (sha256 of the raw file) and `.CustomValues` (whatever the chart passed), plus two template functions that resolve images from `versions.yaml` **and register them for version checking** — this is how image pins that only live in a config file become visible to `check-versions`:

- `{{ Image "org/app" }}` → `org/app:<version>` (full reference)
- `{{ ImageTag "org/app" }}` → `<version>` (tag only; panics at synth time if the image is missing from `versions.yaml`)

`check-versions` and `create-prs` therefore cover images referenced this way; a new such image needs a `versions.yaml` entry alongside the template.

### What each chart generates

A typical chart Go file produces a `dist/<name>.k8s.yaml` containing:
1. `Namespace`
2. `HelmRepository` (Flux source CRD) — always created in the **`flux-system`** namespace, not the chart's own
3. `ConfigMap` with Helm values, hashed into the `configMapHash` annotation as above
4. `HelmRelease` referencing the above, in the chart's namespace
5. Supporting CRDs as needed (ClusterIssuer, ExternalSecret, SecretStore, PodMonitor, etc.)

### Chart categories

`charts/{apps,infra,observability,security,storage}/` — one file per namespace, each exporting a `NewXxxChart(*kubehelpers.Builder) *kubehelpers.Chart` constructor. Registration order in `charts/homelab.go` is grouped by category and roughly reflects bootstrap dependency (security → storage → infra → observability → apps), though Flux, not this ordering, governs actual apply order.

## Important files and directories

- **`charts/homelab.go`** — master registry; a slice of constructor functions. Commenting an entry out is how an app is disabled without deleting its definition.
- **`charts/{apps,infra,observability,security,storage}/`** — one `.go` file per deployed app/namespace
- **`internal/kubehelpers/`** — shared builder library (HelmRelease, Ingress, StatefulSet helpers, etc.)
- **`internal/gitea/`** — Gitea API client behind the `create-prs` / `list-prs` / `merge-pr` subcommands
- **`configs/`** — Helm values YAML per release, injected as ConfigMap into each HelmRelease
- **`versions.yaml`** — single source of truth for all Helm chart and container image versions
- **`imports/`** — auto-generated Go CRD bindings. Gitignored; **do not edit or commit.** Regenerate with `mise run import`. Excluded from linting.
- **`dist/`** — generated output. Gitignored on `main`; CI copies to `generated` branch.
- **`policies/`** — OPA/Rego policies for `conftest` validation. Rules are hard `deny`; legitimate exceptions live in `policies/exemptions.yaml`, keyed by resource name, each with a `reason`. Fix the issue, delete the entry — the policy applies again automatically. New apps must either comply (pinned tag, `runAsNonRoot`) or carry a reasoned exemption; CI (lint workflow) fails on any unexempted violation.
- **`cdk8s.yaml`** — cdk8s config declaring CRD imports and `dist` output path. Rewritten by `cdk8s import`

## Never

- Hand-edit `dist/` or the `generated` branch — CI overwrites both from Go source
- Commit `imports/`, the `k8s-home` binary, or the `cdk8s.yaml` rewrite `cdk8s import` produces
- Update to release candidate, alpha, or beta versions

## Adding a new app

See `docs/adding-new-app.md` for worked examples (stateless app, env vars + Vault secrets, Redis sidecar + ServiceMonitor). Short version:

1. Create `charts/<category>/<appname>.go` with a constructor returning `*kubehelpers.Chart`
2. Register the constructor in `charts/homelab.go`
3. Add Helm values to `configs/<appname>.yaml` (if using Helm)
4. Add version entries to `versions.yaml` under `helmcharts:` or `images:`
5. Run `mise run generate-charts` to produce `dist/<appname>.k8s.yaml`

## Version updates

- Use `mise run check-versions` (or `mise run check-versions --filter <substring>`) to find outdated versions
- Open the PR with `./k8s-home create-prs` (see the command table above; `--dry-run` first)
- Never update to release candidates, alpha, or beta versions
- `excalidraw/excalidraw` image uses `latest` tag — `check-versions` cannot see it; restart the pod to update
- Temporal has a bespoke upgrade path: `contrib/upgrade-temporal.sh` and `docs/temporal.md`
- `check-versions` has known blind spots — see the versions.yaml section above; some containers require manual review

## CI/CD

- **Lint** (`.gitea/workflows/lint.yaml`): runs on PRs to `main`, skipped for `versions.yaml`-only changes
- **Deploy** (`.gitea/workflows/deploy.yaml`): on merge to `main`, runs `cdk8s import` + `cdk8s synth`, pushes `dist/` to `generated` branch, then pushes an OCI artifact
- **Check versions** (`.gitea/workflows/check-versions.yaml`): daily cron + manual dispatch; runs `check-versions` then `create-prs`
- **Build CI image** (`.gitea/workflows/build-image.yaml`): on changes to `contrib/build-image/**`; rebuilds and pushes the CI container image — no cdk8s involved
- The three cdk8s workflows above each do `cdk8s import` → `git checkout cdk8s.yaml` → `cdk8s synth`
- CI container image: `registry.mkz.me/mycroft/golang-cdk8s:latest`
- Env vars needed for PR operations: `GITEA_TOKEN` (primary) or `GITHUB_TOKEN` (fallback)

## Conventions

- **Each chart = one Kubernetes namespace.** Namespace name equals chart/release name.
- **Secrets**: Vault + External Secrets Operator. Vault path pattern: `secret/namespaces/<namespace>/<secret-name>`. Use `kubehelpers.CreateSecretStore()` + `kubehelpers.CreateExternalSecret()`. Sealed Secrets is only used for bootstrapping secrets needed before Vault is available.
- **Storage**: Longhorn with `longhorn-crypto-global` storage class (hardcoded in StatefulSet helpers).
- **Ingress**: Traefik with `cert-manager.io/cluster-issuer: letsencrypt-prod`. Three domains are in use:
  - `*.services.mkz.me` — the default for internal services
  - `*.iop.cx` — public-facing hostnames (e.g. `todo`, `wiki`, `n8n`, `s3-api`); grep `charts/` for the current set rather than trusting a list here
  - `*.lan.mkz.me` — LAN-only hosts
- **Multi-host ingress**: a chart may declare several hostnames (`Ingresses: []string{...}` / `NewAppIngresses`). When it does, check whether the app needs the full host list elsewhere in its config — CORS allow-lists in particular. See `charts/apps/vikunja.go`, which derives `VIKUNJA_CORS_ORIGINS` from the ingress list so a new hostname needs only one edit.
- **Linting**: `.golangci.yaml` uses `default: none` with an explicit enable list; `imports/` is excluded.

## See also

- `README.md` — full project documentation
- `docs/adding-new-app.md` — how to add an app, with worked examples
- `docs/vault.md` — Vault setup and unsealing
- `docs/flux.md` — Flux CD installation, upgrade, troubleshooting
- `docs/cnpg.md` — CloudNativePG (PostgreSQL) operator and clusters
- `docs/garage.md` — Garage S3 storage
- `docs/authentik.md` — Authentik SSO
- `docs/temporal.md` — Temporal, including its upgrade procedure

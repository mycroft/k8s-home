# k8s-home

> cdk8s-based GitOps CLI for managing a personal homelab Kubernetes cluster.

This project uses Go code to programmatically define and generate Flux CD manifests for a k3s cluster of three Mini PCs running Fedora Server. Source code in `charts/` is compiled and synthesized into YAML manifests in `dist/`, which are pushed to the `generated` branch on merge. Flux CD polls that branch every 2 minutes and reconciles the cluster state.

## Features

- **Chart generation** — Go-based cdk8s code synthesizes the manifests of more than 70 apps and infrastructure, observability, security and storage components
- **Version management** — a single `versions.yaml` file tracks all Helm chart and container image versions, with optional regex filters
- **Automated version bumps** — CI checks for outdated versions twice a day and opens a pull request on Gitea for each one
- **GitOps deployment** — a Gitea Actions pipeline builds the charts on merge to `main`; Flux CD applies them to the cluster
- **Policy checks** — OPA/Rego policies and a Flux-style kustomize build validate every pull request
- **Managed PostgreSQL** — CloudNativePG provisions application databases, roles, extensions, and Vault credentials
- **Layered backups** — Longhorn snapshots and backups, Velero/Kopia filesystem backups, and daily logical PostgreSQL dumps
- **Tailnet access** — the Tailscale operator serves selected services on the tailnet, next to their regular ingress; Vaultwarden is moving to tailnet-only access

## Getting Started

### Prerequisites

- [mise](https://mise.jdx.dev/) — installs Go 1.27.2, golangci-lint, node and the cdk8s CLI (declared in `.mise.toml`, pinned by `mise.lock`)
- [conftest](https://www.conftest.dev/) and [kustomize](https://kustomize.io/) (or `kubectl`) — required by `mise run conftest` and `mise run kustomize-build`, installed outside of mise
- `kubectl` configured against the k3s cluster
- `GITEA_TOKEN` (or `GITHUB_TOKEN`) — only for the pull request subcommands
- (Optional, only for the legacy `contrib/create-pr.sh`) `tea` and `just`

### First Run

`imports/` (the Go CRD bindings) is gitignored and nothing builds without it, so generate it first:

```sh
git clone <repo-url>
cd k8s-home
mise run import
mise run generate-charts
```

## Build & Run

| Command                    | What it does                                                                                                  |
| -------------------------- | ------------------------------------------------------------------------------------------------------------- |
| `mise run build`           | Compiles the `k8s-home` binary                                                                                |
| `mise run generate-charts` | Builds, then synthesizes `dist/<chart>.k8s.yaml`                                                              |
| `mise run synth`           | Same output through the cdk8s CLI (what CI runs)                                                              |
| `mise run lint`            | Runs golangci-lint                                                                                            |
| `mise run conftest`        | Generates the charts, then runs the OPA policies in `policies/` against `dist/`                               |
| `mise run kustomize-build` | Generates the charts, then builds `dist/` the way Flux's kustomize-controller does                            |
| `mise run check-versions`  | Lists outdated Helm charts and images                                                                         |
| `mise run import`          | Regenerates `imports/` from the `imports:` list in `cdk8s.yaml`                                               |
| `go test ./...`            | Runs the test suite                                                                                           |
| `mise run diff`            | Diffs `dist/` against the `generated` branch, in a scratch directory (run `generate-charts` first)            |
| `sh contrib/diff.sh`       | Diffs `dist/` against the `generated` branch (checks `generated/` out into the working tree, then deletes it) |

Before committing a Go change:

```sh
mise run lint && mise run generate-charts && mise run conftest && mise run kustomize-build && go test ./...
```

Then inspect the `dist/` diff: CI does not gate on it.

## Version Updates

```sh
mise run check-versions                  # name;old;new for each outdated chart or image
./k8s-home create-prs --dry-run          # the pull requests that would be created
./k8s-home create-prs --filter <name>    # create them on Gitea
./k8s-home list-prs
./k8s-home merge-pr <number>             # rebase-merges: Flux deploys it
```

CI runs `check-versions` then `create-prs` twice a day, so most bumps arrive as pull requests on their own. Changes to `versions.yaml` alone skip the lint workflow: a bad bump is only caught by local verification.

Never update to release candidates, alpha, or beta versions. Some images are not covered by `check-versions`:

- `excalidraw/excalidraw` only publishes `latest`: restart the pod to update it, and check the latest pushed image on [Docker Hub](https://hub.docker.com/r/excalidraw/excalidraw/tags).
- `ghcr.io/karakeep-app/karakeep-chrome` uses a non-semver tag (`154.0.8037.57-r1`): bump it by hand.
- Temporal has its own upgrade path: see [docs/temporal.md](docs/temporal.md).

The older `contrib/create-pr.sh` (`-f <filter>`, `-m 'chart;old;new'`) still works, but discards uncommitted edits to `versions.yaml`, force-pushes, and with `-M` merges straight to `main`.

## Configuration

| Variable       | Required           | Default | Description                                                                           |
| -------------- | ------------------ | ------- | ------------------------------------------------------------------------------------- |
| `GITEA_TOKEN`  | For PR subcommands | —       | Authentication token for the Gitea API (used by `create-prs`, `list-prs`, `merge-pr`) |
| `GITHUB_TOKEN` | No                 | —       | Fallback authentication token if `GITEA_TOKEN` is not set                             |
| `CDK8S_OUTDIR` | No                 | `dist`  | Output directory of the synthesized manifests                                         |

### CLI Flags

| Flag              | Default              | Description                                                                      |
| ----------------- | -------------------- | -------------------------------------------------------------------------------- |
| `--versions-file` | `versions.yaml`      | Path to the versions file                                                        |
| `--debug`         | `false`              | Enable debug logging                                                             |
| `--gitea-url`     | `https://git.mkz.me` | Gitea instance URL                                                               |
| `--owner`         | `mycroft`            | Repository owner                                                                 |
| `--repo`          | `k8s-home`           | Repository name                                                                  |
| `--filter`        | *(empty)*            | Substring matched against chart and image names (`check-versions`, `create-prs`) |
| `--dry-run`       | `false`              | `create-prs`: print the PRs that would be created without creating them          |
| `--branch`        | *(empty)*            | `create-prs`: override the generated branch name (requires exactly one update)   |
| `--base-branch`   | `main`               | `create-prs`: base branch for the PRs                                            |

`--versions-file` with `CDK8S_OUTDIR` synthesizes a speculative bump elsewhere without touching `versions.yaml`: diff the two trees.

## Project Structure

| Path                    | Description                                                                                                                       |
| ----------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| `charts/homelab.go`     | The list of deployed charts; commenting an entry out disables an app                                                              |
| `charts/apps/`          | User-facing app charts (wallabag, freshrss, vaultwarden, etc.)                                                                    |
| `charts/infra/`         | Infrastructure charts (Flux CD, Traefik, Tailscale, Velero, Temporal, etc.)                                                       |
| `charts/observability/` | Monitoring charts (Prometheus, Grafana, Loki, Tempo, etc.)                                                                        |
| `charts/security/`      | Security charts (Vault, cert-manager, Dex, Authentik, Falco, etc.)                                                                |
| `charts/storage/`       | Storage charts (Longhorn, CloudNativePG, NATS, Garage, etc.)                                                                      |
| `internal/kubehelpers/` | Shared builder library for HelmRelease, Ingress, StatefulSet, etc.                                                                |
| `internal/gitea/`       | Gitea API client for PR automation                                                                                                |
| `configs/`              | Helm values per release (flat `<release>.yaml`), injected as ConfigMaps; `configs/<app>/` holds an app's own native config files |
| `versions.yaml`         | Single source of truth for all versions                                                                                           |
| `policies/`             | OPA/Rego policy rules for `conftest`, and their reasoned exemptions                                                               |
| `docs/`                 | Operational guides (see [References](#references))                                                                                |
| `contrib/`              | Helper scripts                                                                                                                    |
| `crds/`                 | Custom CRD YAML files for cdk8s imports                                                                                           |
| `imports/`              | Generated Go CRD bindings (gitignored)                                                                                            |
| `dist/`                 | Generated output, `.k8s.yaml` files (gitignored)                                                                                  |
| `.gitea/workflows/`     | CI: lint, deploy, check-versions, CI image build                                                                                  |

## Installed Components

### Security

- [cert-manager](https://cert-manager.io/) — Let's Encrypt certificates for ingresses
- [Vault](https://www.vaultproject.io/) & [External Secrets Operator](https://external-secrets.io/latest/) — secrets management; [sealed-secrets](https://github.com/bitnami-labs/sealed-secrets) only for secrets needed before Vault is available
- [dex-idp](https://dexidp.io/) — OIDC provider backed by the personal [Gitea](https://about.gitea.com/) instance
- [authentik](https://goauthentik.io/) — SSO for apps, and proxy outposts for admin UIs (see [docs/authentik.md](docs/authentik.md))
- [traefik-forward-auth](https://github.com/thomseddon/traefik-forward-auth) — OIDC forward authentication for admin UIs without their own login
- [Trivy operator](https://github.com/aquasecurity/trivy-operator) — vulnerability reports, browsed through [trivy-operator-dashboard](https://github.com/raoulx24/trivy-operator-dashboard); [harbor-scanner-trivy](https://github.com/goharbor/harbor-scanner-trivy) is the Trivy scanner of the Harbor image registry
- [Falco](https://falco.org/) — runtime syscall detection; its rules are still being tuned, nothing alerts on them yet

### Storage

- [Longhorn](https://longhorn.io/) — distributed, encrypted block storage across the nodes
- [CloudNativePG](https://cloudnative-pg.io/) — PostgreSQL operator and the shared application database cluster (see [docs/cnpg.md](docs/cnpg.md))
- [Garage](https://garagehq.deuxfleurs.fr/) — S3-compatible object storage: the backup bucket lives on a single-node Garage on moonstone, outside the cluster, whose `s3.lan.mkz.me` certificate the cluster issues; `charts/storage/garage.go` still deploys the earlier in-cluster release behind `s3-api.iop.cx` (see [docs/garage.md](docs/garage.md))
- [NATS](https://nats.io/) — message queues

### Infrastructure

- [Flux CD](https://fluxcd.io/) — GitOps reconciliation (see [docs/flux.md](docs/flux.md)), with [fluxcd-ui](https://git.mkz.me/mycroft/fluxcd-ui) as its web UI
- [Traefik](https://traefik.io/) — the ingress controller bundled with k3s; this repository manages its dashboard route
- [Tailscale operator](https://tailscale.com/docs/kubernetes-operator) — serves Vault, Vaultwarden, Atuin, Atuin AI and LiteLLM at `https://<name>.merino-tiaki.ts.net` on the tailnet (see [docs/tailscale.md](docs/tailscale.md))
- [Velero](https://velero.io/) — cluster backup and restore
- [Temporal](https://temporal.io/) — durable workflow engine (see [docs/temporal.md](docs/temporal.md))

### Observability

- [kube-prometheus-stack](https://github.com/prometheus-operator/kube-prometheus) — Prometheus, Alertmanager and [Grafana](https://grafana.com/grafana/)
- [Loki](https://grafana.com/oss/loki/) — log aggregation and storage
- [Tempo](https://grafana.com/oss/tempo/) — trace storage
- [Grafana Alloy](https://grafana.com/docs/alloy/) — collects pod logs and ships them to Loki
- [blackbox-exporter](https://github.com/prometheus/blackbox_exporter) and [smokeping-prober](https://github.com/SuperQ/smokeping_prober) — HTTP and ICMP probing
- [Karma](https://github.com/prymitive/karma) — Alertmanager dashboard
- Custom monitoring — scrape targets outside the cluster (LAN hosts, Garage, the LLM API)

### User Apps

- [Atuin](https://atuin.sh/) — shell history sync server
- [Atuin AI](https://github.com/atuinsh/atuin-ai-server) — backend of `atuin ai`, through LiteLLM
- [Calibre-Web](https://github.com/janeczku/calibre-web) — e-book library manager
- [CyberChef](https://gchq.github.io/CyberChef/) — cyber toolkit
- [Excalidraw](https://excalidraw.com/) — virtual collaborative whiteboard
- [FreshRSS](https://freshrss.org/) — RSS aggregator
- [Homepage](https://github.com/gethomepage/homepage) — dashboard
- [IT-Tools](https://it-tools.tech/) — handy tools for engineers
- [Karakeep](https://karakeep.app/) (formerly Hoarder) — bookmark manager with AI tagging
- [Linkding](https://github.com/sissbruecker/linkding) — bookmark manager
- [LiteLLM](https://docs.litellm.ai/) — OpenAI-compatible proxy in front of llama-server on moonstone and OpenAI
- [Memos](https://www.usememos.com/) — note taking
- [MicroBin](https://microbin.eu/) — lightweight pastebin
- [n8n](https://n8n.io/) — workflow automation
- [Open WebUI](https://openwebui.com/) — AI chat interface
- [Opengist](https://github.com/thomiceli/opengist) — GitHub Gist clone
- [Outline](https://www.getoutline.com/) — team knowledge base
- [Paperless-ngx](https://docs.paperless-ngx.com/) — document management
- [PrivateBin](https://privatebin.info/) — secure pastebin
- [Send](https://gitlab.com/timvisee/send) — simple, private file sharing
- [Vaultwarden](https://github.com/dani-garcia/vaultwarden) — Bitwarden-compatible password manager, moving to tailnet-only access
- [Vikunja](https://vikunja.io/) — task management
- [Wallabag](https://www.wallabag.it/) — read-it-later list
- [Wiki.js](https://js.wiki/) — wiki platform
- [Yopass](https://yopass.se/) — secure secret sharing
- [Zipline](https://github.com/diced/zipline) — file sharing service

### Agents

- [agent-nexus](https://git.mkz.me/mycroft/agent-nexus) — LLM agent that reviews pull requests
- [Zulip](https://zulip.com/) — team chat used to talk to agent-nexus
- [reviews](https://git.mkz.me/mycroft/reviews) — dashboard of the GitHub pull requests waiting for review
- Read-only MCP servers agent-nexus inspects the cluster through: [flux-operator-mcp](https://github.com/controlplaneio-fluxcd/flux-operator), [kubernetes-mcp-server](https://github.com/containers/kubernetes-mcp-server) and [mcp-grafana](https://github.com/grafana/mcp-grafana)

### Toys, Demos & Redirections

- Llama — ExternalName ingress redirecting `llm-api.iop.cx` to a self-hosted LLM API
- Music Assistant — ExternalName ingress redirecting `music-assistant.iop.cx` to a self-hosted server
- Useless — collection of useless web pages
- Root Me — toy scoreboard at `root-me.iop.cx`
- URLs — curated list of useful links
- Hey — self-built hello page at `hey.iop.cx`
- app-sample and hello-kubernetes — demo workloads

### Disabled

Defined in `charts/` but commented out in `charts/homelab.go`: Kyverno, MinIO, ScyllaDB, MariaDB, OpenSearch, Capacitor, Linkerd, Promtail, Jaeger, BookStack, Emojivoto, Redmine, Snippet Box and Heimdall.

## Adding an App

See [docs/adding-new-app.md](docs/adding-new-app.md) for worked examples. In short:

1. Create `charts/<category>/<app>.go` with a constructor returning `*kubehelpers.Chart`
2. Register the constructor in `charts/homelab.go`
3. Add Helm values to `configs/<app>.yaml` (if using Helm)
4. Add the version entries to `versions.yaml`
5. Run `mise run generate-charts` and the checks above

## Deployment

- **Lint** (`.gitea/workflows/lint.yaml`) — on pull requests to `main`, except `versions.yaml`-only changes: synthesizes the charts, then runs conftest, the kustomize build, golangci-lint and `go test ./...`.
- **Deploy** (`.gitea/workflows/deploy.yaml`) — on merge to `main`: synthesizes the charts and replaces the content of the `generated` branch with `dist/`. Flux polls that branch every 2 minutes. Its Kustomization prunes, so removing a chart deletes its namespace and everything in it. The workflow also publishes an OCI artifact, which no Kustomization consumes yet (see [docs/flux.md](docs/flux.md)).
- **Check versions** (`.gitea/workflows/check-versions.yaml`) — twice a day: opens the version-bump pull requests.
- **Build CI image** (`.gitea/workflows/build-image.yaml`) — on changes to `contrib/build-image/`: rebuilds the CI container image. The other workflows pin it by digest, so a rebuild needs those pins, and `mise.lock`, updated.

## Security

- **Secrets** live in Vault at `secret/namespaces/<namespace>/<name>` and reach the cluster as Kubernetes `Secret`s through External Secrets. `sealed-secrets` only bootstraps what is needed before Vault is available.
- **Access**: apps with OIDC support log in through authentik or dex. Admin UIs without their own login sit behind `traefik-forward-auth` or an authentik proxy outpost. A few services are also served on the tailnet; Vaultwarden's public ingress goes once every client has moved there (see [docs/tailscale.md](docs/tailscale.md#vaultwarden-migration)).
- **Policies**: `conftest` rejects, among others, unpinned or `latest` images, root or privileged containers, ingresses without TLS from `letsencrypt-prod` or outside the fleet domains, `NodePort` and `LoadBalancer` services, unencrypted storage classes, and inline secrets. Legitimate exceptions are listed, each with a reason, in `policies/exemptions.yaml`.
- **Runtime**: Falco watches syscalls on every node; the Trivy operator reports vulnerabilities in running images.

## Services

### Traefik Configuration

Traefik is installed by k3s. Its values are set by hand on the server node, in `/var/lib/rancher/k3s/server/manifests/traefik-config.yaml`; k3s redeploys the chart when the file changes. Do not edit `traefik.yaml` directly — it is overwritten on restart. The current values:

```yaml
apiVersion: helm.cattle.io/v1
kind: HelmChartConfig
metadata:
  name: traefik
  namespace: kube-system
spec:
  valuesContent: |-
    logs:
      access:
        enabled: true
    additionalArguments:
      - "--serverstransport.insecureskipverify=true"
    providers:
      kubernetesCRD:
        allowCrossNamespace: true       # IngressRoute middlewares from other namespaces
        allowExternalNameServices: true # ExternalName services (Llama, Music Assistant)
    ports:
      opengist-ssh:
        port: 22222
        export: true
        protocol: TCP
        exposedPort: 22222
        expose:
          default: true
      websecure:
        http3:
          enabled: true
      web:
        redirections:
          entryPoint:
            to: websecure
            scheme: https
            permanent: true
```

See the [Traefik Helm chart values](https://github.com/traefik/traefik-helm-chart/blob/master/traefik/values.yaml) for more options.

### Traefik Dashboard

The dashboard is at `https://traefik.services.mkz.me/dashboard/`, behind `traefik-forward-auth`. Its `IngressRoute` and certificate are defined in `charts/infra/traefik.go`.

### Vault and External-Secrets

Vault is configured for Kubernetes-based authentication with a custom policy for External-Secrets queries. Transient service account tokens eliminate the need to maintain service account secrets. See the [External-Secrets Vault provider docs](https://external-secrets.io/latest/provider/hashicorp-vault/) and [Vault Kubernetes auth docs](https://developer.hashicorp.com/vault/docs/auth/kubernetes).

Setup steps:

```sh
vault auth enable kubernetes

vault write auth/kubernetes/config \
    kubernetes_host=https://10.43.0.1:443

vault policy write external-secrets contrib/external-secrets.vault.hcl

vault write auth/kubernetes/role/external-secrets \
    bound_service_account_names=external-secrets \
    bound_service_account_namespaces=external-secrets \
    policies=external-secrets \
    ttl=24h
```

This is sufficient to create valid `SecretStore` and `ExternalSecret` objects. [docs/vault.md](docs/vault.md) has the full procedure, including initialization, unsealing and OIDC login.

Vault is served at `https://vault.services.mkz.me` and, on the tailnet, at `https://vault.merino-tiaki.ts.net`. The CLI's OIDC login works at both; the UI's needs an extra redirect URI for the tailnet address (see [docs/vault.md](docs/vault.md)).

## Operational Notes

### Creating a Secret with Sealed-Secrets

See [bitnami-labs/sealed-secrets](https://github.com/bitnami-labs/sealed-secrets):

```sh
echo -n bar | kubectl create secret generic mysecret --dry-run=client --from-file=foo=/dev/stdin -o json > mysecret.json
kubeseal < mysecret.json > mysealedsecret.json
kubectl create -f mysealedsecret.json
```

Warning: Do not forget the namespace and secret name. SealedSecret is strict about both.

### Creating PostgreSQL Storage

1. Add a `CNPGDatabase` entry to `charts/storage/cnpg-cluster.go`.
2. Add an `ExternalSecret` in the application chart for the configured Vault entry.
3. Commit and push. CloudNativePG creates the database and role, while External Secrets publishes the generated credentials to Vault.
4. Read the synchronized Secret from the application workload. It provides `host`, `dbname`, `username`, `password`, and `url` keys.

See [docs/cnpg.md](docs/cnpg.md) for field definitions, credential flow,
backup verification, and restore guidance.

### Checking Backups

PostgreSQL logical dumps run at `00:30` and the cluster-wide Velero/Kopia
filesystem backup runs at `01:30`; a second Velero schedule backs up the
cluster resources at `07:30`. Both Velero backups are kept for 30 days.
Longhorn also takes snapshots four times a day and a backup to S3 at `03:30`.
Kubernetes CronJobs use UTC unless a `timeZone` is configured.

```sh
kubectl get cronjob,job -n cnpg -l app.kubernetes.io/name=postgres-dump
kubectl get backups.velero.io -n velero --sort-by=.metadata.creationTimestamp
```

While the dump pod is active, the Velero filesystem backup validates that its
logical dump is complete and fresh before copying the PVC. Check both the dump
Job and Velero backup; see [docs/cnpg.md](docs/cnpg.md#backups) for the manual
verification and recovery procedure.

## References

- [cdk8s Documentation](https://cdk8s.io/) — Cloud Development Kit for Kubernetes
- [Flux CD Documentation](https://fluxcd.io/) — GitOps Kubernetes automation
- [External-Secrets Vault Provider](https://external-secrets.io/latest/provider/hashicorp-vault/) — Vault integration guide
- [Vault Kubernetes Auth](https://developer.hashicorp.com/vault/docs/auth/kubernetes) — Kubernetes authentication for Vault
- [Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets) — Encrypted Kubernetes secrets
- [Traefik Helm Chart Values](https://github.com/traefik/traefik-helm-chart/blob/master/traefik/values.yaml) — Traefik configuration reference
- [docs/adding-new-app.md](docs/adding-new-app.md) — adding an app, with worked examples
- [docs/flux.md](docs/flux.md) — Flux CD installation, upgrade, and troubleshooting playbook
- [docs/vault.md](docs/vault.md) — Vault setup, initialization, and unsealing
- [docs/authentik.md](docs/authentik.md) — authentik SSO and proxy outposts
- [docs/cnpg.md](docs/cnpg.md) — CloudNativePG databases and their credentials
- [docs/garage.md](docs/garage.md) — Garage S3 storage
- [docs/temporal.md](docs/temporal.md) — Temporal workflow setup and upgrade
- [docs/tailscale.md](docs/tailscale.md) — Tailscale operator setup and tailnet-only services

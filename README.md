# k8s-home

> cdk8s-based GitOps CLI for managing a personal homelab Kubernetes cluster.

This project uses Go code to programmatically define and generate Flux CD `HelmRelease` manifests for a k3s cluster running on multiple Mini PCs with Fedora. Source code in `charts/` is compiled and synthesized into YAML manifests in `dist/`, which are pushed to the `generated` branch on merge. Flux CD polls that branch every 2 minutes and reconciles the cluster state.

## Features

- **Chart generation** — Go-based cdk8s code synthesizes Kubernetes manifests for 35+ apps, infrastructure, observability, and security components
- **Version management** — Single `versions.yaml` file tracks all Helm chart and container image versions with optional regex filters
- **Automated PR creation** — CLI commands check for outdated versions and create pull requests on Gitea
- **GitOps deployment** — Gitea Actions pipeline builds charts on merge to `main`, Flux CD applies them to the cluster

## Getting Started

### Prerequisites

- Go 1.23+ (toolchain go1.24.2)
- [just](https://just.systems/) command runner
- [golangci-lint](https://golangci-lint.run/) for linting
- [cdk8s](https://cdk8s.io/) CLI for importing CRDs
- kubectl configured against the target k3s cluster
- Environment variables: `GITEA_TOKEN` or `GITHUB_TOKEN`

### Installation

```sh
git clone <repo-url>
cd k8s-home
```

```sh
export GITEA_TOKEN=<your-gitea-token>
```

## Build & Run

### Generate Charts

```sh
just generate
```

This builds the binary and runs `./k8s-home` to synthesize `dist/*.k8s.yaml` files.

### Production Build

```sh
just build
```

### Lint

```sh
just lint
```

### Check for Outdated Versions

```sh
just check-versions
```

### Create Update PR

```sh
just update-version <project/component>
```

### Diff Generated Output

```sh
just diff
```

Generates charts then diffs `dist/` against the `generated` branch.

### Import CRDs

```sh
just import
```

Regenerates Go CRD bindings in `imports/`.

### Tests

```sh
go test ./...
```

## Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `GITEA_TOKEN` | Yes | — | Authentication token for Gitea API (used by `create-prs`, `list-prs`, `merge-pr`) |
| `GITHUB_TOKEN` | No | — | Fallback authentication token if `GITEA_TOKEN` is not set |

### CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--versions` | `versions.yaml` | Path to the versions file |
| `--debug` | `false` | Enable debug logging |
| `--gitea-url` | `https://git.mkz.me` | Gitea instance URL |
| `--owner` | `mycroft` | Repository owner |
| `--repo` | `k8s-home` | Repository name |

## Project Structure

| Path | Description |
|------|-------------|
| `charts/apps/` | User-facing app charts (wallabag, freshrss, vaultwarden, etc.) |
| `charts/infra/` | Infrastructure charts (Flux CD, Traefik, Velero, Temporal, etc.) |
| `charts/observability/` | Monitoring charts (Grafana, Prometheus, Loki, Tempo, etc.) |
| `charts/security/` | Security charts (Vault, cert-manager, Dex, Authentik, etc.) |
| `charts/storage/` | Storage charts (Longhorn, PostgreSQL, NATS, Garage, etc.) |
| `charts/static/` | Static YAML manifests (Tekton pipeline definitions) |
| `internal/kubehelpers/` | Shared builder library for HelmRelease, Ingress, StatefulSet, etc. |
| `internal/gitea/` | Gitea API client for PR automation |
| `configs/` | Helm values YAML files, injected as ConfigMaps |
| `dist/` | Generated output (`.k8s.yaml` files) |
| `imports/` | Auto-generated Go CRD bindings |
| `crds/` | Custom CRD YAML files for cdk8s imports |
| `cicd/` | Tekton pipeline definitions |
| `contrib/` | Helper scripts |
| `versions.yaml` | Single source of truth for all versions |

## Installed Apps

Most important apps installed on the cluster:

- [longhorn](https://longhorn.io/) — distributed block storage across the cluster
- [vault](https://www.vaultproject.io/) & [External Secrets Operator](https://external-secrets.io/latest/) — secrets management; [sealed-secrets](https://github.com/bitnami-labs/sealed-secrets) for encrypted volume bootstrapping
- [dex-idp](https://dexidp.io/) — OAuth SSO linked with personal [gitea](https://about.gitea.com/) instance
- [authentik](https://goauthentik.io/) — SSO authentication and authorization
- [traefik-forward-auth](https://doc.traefik.io/traefik/middlewares/http/forwardauth/) — OAuth authn/authz firewall for apps not linked to dex-idp
- [cert-manager](https://cert-manager.io/) — on-demand TLS certificate generation for ingresses
- [PostgreSQL](https://www.postgresql.org/) operator — database instances
- [NATS](https://nats.io/) — message queues
- [kube-prometheus-stack](https://github.com/prometheus-operator/kube-prometheus) — metrics monitoring with [Grafana](https://grafana.com/grafana/)
- [blackbox-exporter](https://github.com/prometheus/blackbox_exporter) — blackbox probing
- Grafana [loki](https://grafana.com/oss/loki/) — log aggregation and storage
- [tempo](https://grafana.com/oss/tempo/) — trace processing
- [Grafana Alloy](https://grafana.com/docs/alloy/) — observability agent
- [Karma](https://github.com/prymitive/karma) — Prometheus alert dashboard
- [Garage](https://garagehq.deuxfleurs.fr/) — S3-compatible object storage
- [Velero](https://velero.io/) — cluster backup and restore
- [Capacitor](https://capacitor.l5d.io/) — in-cluster CI/CD

### User Apps

- [wallabag](https://www.wallabag.it/) — visual to-read list
- [freshrss](https://freshrss.org/) — RSS aggregator
- [privatebin](https://privatebin.info/) — secure pastebin
- [paperless-ngx](https://docs.paperless-ngx.com/) — document management
- [yopass](https://yopass.se/) — secure secret sharing
- [bookstack](https://www.bookstackapp.com/) — information organization platform
- [IT-Tools](https://it-tools.tech/) — handy tools for engineers
- [vaultwarden](https://github.com/dani-garcia/vaultwarden) — Bitwarden-compatible password manager
- [send](https://gitlab.com/timvisee/send) — simple, private file sharing
- [snippetbox](https://github.com/pawelmalak/snippet-box) — code snippet portal
- [excalidraw](https://excalidraw.com/) — virtual collaborative whiteboard
- [wikijs](https://js.wiki/) — wiki platform
- [redmine](https://www.redmine.org/) — project management
- [microbin](https://microbin.eu/) — lightweight pastebin
- [memos](https://www.usememos.com/) — note taking
- [opengist](https://github.com/thomiceli/opengist) — GitHub Gist clone
- [Hoarder](https://github.com/hoarder-app/hoarder) — bookmark manager
- [Vikunja](https://vikunja.io/) — task management
- [Open WebUI](https://open-webui.com/) — AI chat interface
- [Zipline](https://github.com/diced/zipline) — file sharing service
- [Calibre Web](https://calibre-ebook.com/download) — e-book library manager
- [CyberChef](https://gchq.github.io/CyberChef/) — cyber toolkit
- [Homepage](https://github.com/gethomepage/homepage) — dashboard
- [Happy-Urls](https://github.com/awesome-selfhosted/happy-urls) — URL shortener
- [Linkding](https://github.com/sissbruecker/linkding) — bookmark manager
- [n8n](https://n8n.io/) — workflow automation

## Deployment

When a PR is merged into `main`, a Gitea Actions pipeline builds cdk8s dependencies, generates charts, and pushes them to the `generated` branch. Flux CD pulls the branch every 2 minutes. See `.gitea/workflows/deploy.yaml`.

## Maintenance

### Upgrade Services

Check for outdated Helm charts and images:

```sh
just check-versions
```

Run updates:

```sh
just update-version <project/component>
```

`<project/component>` must include the full project name and component name. Never update to release candidates, alpha, or beta versions. Note that not all containers are checked by this command.

#### Excalidraw

The `excalidraw/excalidraw` image has no tag. Restart the pod to get updates. Check the latest pushed image on [Docker Hub](https://hub.docker.com/r/excalidraw/excalidraw/tags).

## Security

Secrets are stored either in source code via `sealed-secrets` or in Vault and deployed as Kubernetes `Secret` resources through `external-secrets`.

Ingresses are protected by `traefik-forward-auth`, which performs OAuth using a private and external Gitea instance.

## Services

### Traefik Configuration

Redirect all HTTP requests to HTTPS. On the master node, edit `/var/lib/rancher/k3s/server/manifests/traefik-config.yaml` and add in the `valuesContent` section:

```yaml
  valuesContent: |-
    ports:
      web:
        redirectTo: websecure
```

k3s will automatically redeploy the Helm chart.

Allow `IngressRoute` middlewares from other namespaces:

```yaml
    providers:
      kubernetesCRD:
        allowCrossNamespace: true
```

Do not edit `traefik.yaml` directly — it is overwritten on restart.

### Viewing the Traefik Dashboard

The dashboard is available at `https://traefik.services.mkz.me/dashboard/#/`. Add this `IngressRoute` first:

```yaml
apiVersion: traefik.io/v1alpha1
kind: IngressRoute
metadata:
  name: traefik-dashboard
spec:
  routes:
    - match: Host(`traefik.services.mkz.me`) && (PathPrefix(`/api`) || PathPrefix(`/dashboard`))
      kind: Rule
      services:
        - name: api@internal
          kind: TraefikService
      middlewares:
        - name: traefik-forward-auth-traefik-forward-auth@kubernetescrd
```

### Vault and External-Secrets

Vault is configured for Kubernetes-based authentication with a custom policy for External-Secrets queries. Transient service account tokens eliminate the need to maintain service account secrets. See the [External-Secrets Vault provider docs](https://external-secrets.io/latest/provider/hashicorp-vault/) and [Vault Kubernetes auth docs](https://developer.hashicorp.com/vault/docs/auth/kubernetes).

Setup steps:

```sh
vault auth enable kubernetes

vault write auth/kubernetes/config \
    kubernetes_host=https://10.43.0.1:443

vault write external-secrets readwrite /path/to/external-secrets.vault.hcl

vault write auth/kubernetes/role/external-secrets \
    bound_service_account_names=external-secrets \
    bound_service_account_namespaces=external-secrets \
    policies=external-secrets \
    ttl=24h
```

This is sufficient to create valid `SecretStore` and `ExternalSecret` objects.

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

1. In `postgresql.go`, add a record to the list. Commit and push to spawn the database and create secrets in the `postgres` namespace.
2. When the database is ready, connect with psql to grant table creation:

```sh
k exec -ti -n postgres postgres-instance-0 -- psql -U dex-admin dex
dex=# GRANT CREATE ON SCHEMA public TO PUBLIC;
```

3. Retrieve credentials from the `Secret` in the `postgres` namespace:

```sh
k get secret -n postgres dex-admin.postgres-instance.credentials.postgresql.acid.zalan.do -o yaml | yq -r .data.username | base64 -d
k get secret -n postgres dex-admin.postgres-instance.credentials.postgresql.acid.zalan.do -o yaml | yq -r .data.password | base64 -d
```

4. Create a record in Vault (e.g., `namespaces/<namespace>/postgresql`) with `username` and `password` keys.
5. In the target namespace, create an `ExternalSecretStore` and an `ExternalSecret`. Examples exist throughout the repository.

### Install whatismyip

Disable the default `externalTrafficPolicy` to get real IP:

```yaml
apiVersion: helm.cattle.io/v1
kind: HelmChartConfig
metadata:
  name: traefik
  namespace: kube-system
spec:
  valuesContent: |-
    additionalArguments:
      - "--serverstransport.insecureskipverify=true"
    service:
      spec:
        externalTrafficPolicy: Local
```

See [traefik values.yaml](https://github.com/traefik/traefik-helm-chart/blob/master/traefik/values.yaml) for more options.

## References

- [cdk8s Documentation](https://cdk8s.io/) — Cloud Development Kit for Kubernetes
- [Flux CD Documentation](https://fluxcd.io/) — GitOps Kubernetes automation
- [docs/flux.md](docs/flux.md) — Flux CD installation, upgrade, and troubleshooting playbook
- [External-Secrets Vault Provider](https://external-secrets.io/latest/provider/hashicorp-vault/) — Vault integration guide
- [Vault Kubernetes Auth](https://developer.hashicorp.com/vault/docs/auth/kubernetes) — Kubernetes authentication for Vault
- [Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets) — Encrypted Kubernetes secrets
- [Traefik Helm Chart Values](https://github.com/traefik/traefik-helm-chart/blob/master/traefik/values.yaml) — Traefik configuration reference
- [CLAUDE.md](CLAUDE.md) — Architecture, conventions, and how to add new apps
- [README.vault.md](README.vault.md) — Vault setup, initialization, and unsealing
- [README.cdk8s.md](README.cdk8s.md) — cdk8s-specific documentation
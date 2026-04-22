# Flux CD Playbook

This playbook covers Flux CD installation, upgrade, and troubleshooting for the k8s-home cluster.

## Overview

Flux CD is the GitOps engine that reconciles the cluster state with the manifests in the `generated` branch. The workflow:

1. Go code in `charts/` synthesizes YAML manifests into `dist/`
2. On merge to `main`, Gitea Actions pushes `dist/` to the `generated` branch
3. Flux CD polls the `generated` branch every 2 minutes and applies changes

The `flux-system` namespace contains the Flux controllers. A `PodMonitor` in `charts/infra/fluxcd.go` exposes metrics to Prometheus.

## Prerequisites

- [flux CLI](https://fluxcd.io/flux/install/) installed locally
- kubectl configured against the k3s cluster
- SSH private key file and corresponding public key registered as a deploy key on Gitea
- Gitea repository with read/write access

## Initial Installation

### 1. Install the Flux CLI

Follow the [official getting started guide](https://fluxcd.io/flux/get-started/).

### 2. Bootstrap Flux with Git

Generate an SSH key pair and register the public key as a deploy key on your Gitea repository, then run:

```sh
flux bootstrap git \
  --url=ssh://git@git.mkz.me/mycroft/k8s-home.git \
  --branch=generated \
  --private-key-file=/tmp/rsa \
  --path=generated/
```

This creates the `flux-system` namespace and installs all Flux controllers. The bootstrap command stores the SSH private key as a Kubernetes secret named `flux-system` in the `flux-system` namespace.

### 3. Verify Installation

```sh
flux get all -n flux-system
```

All components should show `revision` and `ready` status.

## Upgrade Flux

### 1. Export the SSH Key

The private key is stored in the `flux-system` secret. Extract it before re-bootstrapping:

```sh
kubectl get secret -n flux-system flux-system -o yaml | yq '.data.identity' -r | base64 -d > /tmp/rsa
```

### 2. Upgrade the Flux CLI

```sh
# macOS
brew upgrade fluxcd/tap/flux

# Linux
sudo snap install flux --classic
# or follow https://fluxcd.io/flux/install/
```

### 3. Migrate Custom Resources

Before re-bootstrapping, migrate Flux custom resources to their latest API version. This is required when upgrading across minor versions.

Migrate resources in the cluster:

```sh
flux migrate
```

Migrate resources in the local repository (the manifests that will be pushed to Git):

```sh
flux migrate -f . --yes
```

### 4. Re-run Bootstrap

```sh
flux bootstrap git \
  --url=ssh://git@git.mkz.me/mycroft/k8s-home.git \
  --branch=generated \
  --private-key-file=/tmp/rsa \
  --path=generated/
```

The bootstrap command upgrades all Flux controllers in-place.

### 5. Verify Upgrade

```sh
flux --version
flux get all -n flux-system
```

## Troubleshooting

### Force HelmRelease Reconcile

After an undetected change, annotate the HelmRelease to trigger reconciliation:

```sh
kubectl annotate --overwrite -n <namespace> helmrelease/<name> reconcile.fluxcd.io/requestedAt="$(date +%s)"
```

### HelmRelease Install Retries Exhausted

After fixing the underlying issue, suspend and resume the HelmRelease:

```sh
flux suspend hr -n <namespace> <name>
flux resume hr -n <namespace> <name>
```

### Check Component Status

```sh
# All components
flux get all -n flux-system

# Specific component
flux get helmrelease -A

# Detailed status
kubectl describe helmrelease -n <namespace> <name>
```

### View Controller Logs

```sh
kubectl logs -n flux-system -l app=helm-controller --tail=100
kubectl logs -n flux-system -l app=source-controller --tail=100
```

### Monitoring

Flux controllers expose Prometheus metrics on the `http-prom` port. The `PodMonitor` defined in `charts/infra/fluxcd.go` scrapes these metrics. Alert rules are configured in `configs/kube-prometheus-stack.yaml` under the `flux-alerts` section.

## References

- [Flux CD Documentation](https://fluxcd.io/)
- [Flux Getting Started](https://fluxcd.io/flux/get-started/)
- [Flux Bootstrap Generic Git](https://fluxcd.io/flux/installation/bootstrap/generic-git-server/)
- [Flux Monitoring Example](https://github.com/fluxcd/flux2-monitoring-example)
- [Flux HelmRelease Troubleshooting](https://fluxcd.io/flux/tutorials/hello-world/)

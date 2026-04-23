# Flux CD Playbook

This playbook covers Flux CD installation, upgrade, and troubleshooting for the k8s-home cluster.

## Overview

Flux CD is the GitOps engine that reconciles the cluster state with manifests. The workflow:

1. Go code in `charts/` synthesizes YAML manifests into `dist/`
2. On merge to `main`, Gitea Actions pushes `dist/` to the `generated` branch and publishes an OCI artifact
3. Flux CD polls the source (Git or OCI) every 2 minutes and applies changes

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

## Migrate from Git to OCI

This procedure switches Flux from syncing the `generated` Git branch to pulling manifests from an OCI registry, with no downtime.

### Prerequisites

- OCI registry credentials (if the registry requires authentication)
- OCI artifacts already pushed to `registry.mkz.me/mycroft/k8s-home/manifests` (tagged `latest`)

### 1. Create Registry Credentials Secret

If your registry requires authentication:

```sh
kubectl create secret docker-registry oci-registry-credentials \
  --docker-server=registry.mkz.me \
  --docker-username=<username> \
  --docker-password=<password> \
  -n flux-system
```

### 2. Create OCIRepository

Apply the OCIRepository resource:

```yaml
apiVersion: source.toolkit.fluxcd.io/v1
kind: OCIRepository
metadata:
  name: k8s-home
  namespace: flux-system
spec:
  interval: 2m0s
  url: oci://registry.mkz.me/mycroft/k8s-home/manifests
  ref:
    tag: latest
  secretRef:
    name: oci-registry-credentials
```

```sh
kubectl apply -f - <<'EOF'
apiVersion: source.toolkit.fluxcd.io/v1
kind: OCIRepository
metadata:
  name: k8s-home
  namespace: flux-system
spec:
  interval: 2m0s
  url: oci://registry.mkz.me/mycroft/k8s-home/manifests
  ref:
    tag: latest
  secretRef:
    name: oci-registry-credentials
EOF
```

### 3. Verify OCIRepository is Ready

```sh
flux get ocirepository k8s-home -n flux-system
```

Wait until the status shows `ready` and the revision matches the latest artifact.

### 4. Update the Kustomization

Patch the existing `flux-system` Kustomization to reference the OCIRepository instead of the GitRepository. The path must also be updated since OCI files are stored at the root (not under `generated/`):

```sh
kubectl patch kustomization flux-system -n flux-system --type=merge -p '{
  "spec": {
    "path": ".",
    "sourceRef": {
      "apiVersion": "source.toolkit.fluxcd.io/v1",
      "kind": "OCIRepository",
      "name": "k8s-home"
    }
  }
}'
```

### 5. Verify Sync

```sh
flux get kustomization flux-system -n flux-system
flux get all -n flux-system
```

Confirm the Kustomization is reconciling successfully against the OCIRepository.

### 6. Clean Up (Optional)

Once the OCI-based sync is verified working, remove the old GitRepository:

```sh
kubectl delete gitrepository flux-system -n flux-system
```

### Rollback

If the OCI sync fails, revert the Kustomization to use the GitRepository:

```sh
kubectl patch kustomization flux-system -n flux-system --type=merge -p '{
  "spec": {
    "path": "generated/",
    "sourceRef": {
      "apiVersion": "source.toolkit.fluxcd.io/v1",
      "kind": "GitRepository",
      "name": "flux-system"
    }
  }
}'
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
- [Flux OCI Repositories](https://fluxcd.io/flux/components/source/ocirepositories/)
- [Flux Monitoring Example](https://github.com/fluxcd/flux2-monitoring-example)
- [Flux HelmRelease Troubleshooting](https://fluxcd.io/flux/tutorials/hello-world/)

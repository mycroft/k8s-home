# k8s-home

This repository builds & stores k8s configuration.

## Installation

Please start `fluxcd`:

```sh
flux bootstrap git \
  --url=ssh://git@git.mkz.me/mycroft/k8s-home.git \
  --private-key-file=/tmp/rsa \
  --branch=generated \
  --path=generated/
```

## Notes

### Force a HelmRelease reconcile after an undetected change

```sh
kubectl annotate --overwrite -n monitoring helmrelease/prometheus reconcile.fluxcd.io/requestedAt="$(date +%s)"
```
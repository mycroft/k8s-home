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

### Install whatismyip

It is required to not use the default externalTrafficPolicy to be able to get real-ip.

Please use the following configuration:

```sh
# cat /var/lib/rancher/k3s/server/manifests/traefik-config.yaml
apiVersion: helm.cattle.io/v1
kind: HelmChartConfig
metadata:
  name: traefik
  namespace: kube-system
spec:
  valuesContent: |-
    service:
      spec:
        externalTrafficPolicy: Local

# kubectl apply -f /var/lib/rancher/k3s/server/manifests/traefik-config.yaml
# kubectl get svc -o yaml -n kube-system traefik | grep -i policy
  externalTrafficPolicy: Local
  internalTrafficPolicy: Cluster
```
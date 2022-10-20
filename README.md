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

Add/edit the annotation:

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

### Creating a secret using Sealed-Secrets

See https://github.com/bitnami-labs/sealed-secrets:

```sh
# Create a json/yaml-encoded Secret somehow:
# (note use of `--dry-run` - this is just a local file!)
echo -n bar | kubectl create secret generic mysecret --dry-run=client --from-file=foo=/dev/stdin -o json >mysecret.json

# Warning! Do not forget the namespace & the secret name. By default SealedSecret is very strict about those.

# This is the important bit:
# (note default format is json!)
kubeseal <mysecret.json >mysealedsecret.json

# At this point mysealedsecret.json is safe to upload to Github,
# post on Twitter, etc.

# Eventually:
kubectl create -f mysealedsecret.json
```
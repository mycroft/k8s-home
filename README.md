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

## Upgrade flux

Dump key file:

```sh
k get secret -n flux-system -o yaml flux-system | yq .data.identity -r | base64 -d > /tmp/rsa
```

Then re-run the `bootstrap` command.

## Upgrade services

To check outdated Helm charts, use:

```sh
go build && ./k8s-home -check-version
```

## Services

### Traefik tweaks

traefik should/must be reconfigured to redirect all `http://` requests to `https://`.

On master node, edit `/var/lib/rancher/k3s/server/manifests/traefik-config.yaml` and add the following snippet in `valuesContent` section:

```
  valuesContent: |-
    ports:
      web:
        redirectTo: websecure
```

k3s will automatically redeploy the helm chart.


Also, it should allow `IngressRoutes` to use `middlewares` from another namespaces the ingressroutes is in (disabled by default)

```
    providers:
      kubernetesCRD:
        allowCrossNamespace: true
```

Note: Do not edit `traefik.yaml`. It is overwritten on restart.


### Viewing traefik dashboard

port-forward the 9000 port and reach `http://localhost:9000/dashboard/`:

```
> k port-forward -n kube-system (k get pods -n kube-system | grep ^traefik | cut -d' ' -f1) 9000:9000
Forwarding from 127.0.0.1:9000 -> 9000
Forwarding from [::1]:9000 -> 9000
```


### Kubernetes-Dashboard

Generate a login and use it on https://kubernetes-dashboard.services.mkz.me/

```sh
kubectl -n kubernetes-dashboard create token admin
```

You need a recent kubectl binary to be able to do this!

This is fully described on https://github.com/kubernetes/dashboard/blob/master/docs/user/access-control/creating-sample-user.md

See https://github.com/kubernetes/dashboard/blob/master/docs/user/access-control/README.md#login-view

### Other services

- Grafana
- Vault


## Notes

### Force a HelmRelease reconcile after an undetected change

Add/edit the annotation:

```sh
kubectl annotate --overwrite -n monitoring helmrelease/prometheus reconcile.fluxcd.io/requestedAt="$(date +%s)"
```


### On HelmRelease 'install retries exhausted' 

After fixing the issue, suspend & resume the installation:

```sh
> flux suspend hr -n kubernetes-dashboard kubernetes-dashboard
► suspending helmrelease kubernetes-dashboard in kubernetes-dashboard namespace
✔ helmrelease suspended

> flux resume hr -n kubernetes-dashboard kubernetes-dashboard
► resuming helmrelease kubernetes-dashboard in kubernetes-dashboard namespace
✔ helmrelease resumed
◎ waiting for HelmRelease reconciliation
✔ HelmRelease reconciliation completed
✔ applied revision 5.11.0
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
    additionalArguments:
      - "--serverstransport.insecureskipverify=true"
    service:
      spec:
        externalTrafficPolicy: Local

# kubectl apply -f /var/lib/rancher/k3s/server/manifests/traefik-config.yaml
# kubectl get svc -o yaml -n kube-system traefik | grep -i policy
  externalTrafficPolicy: Local
  internalTrafficPolicy: Cluster
```

See https://github.com/traefik/traefik-helm-chart/blob/master/traefik/values.yaml for more values.


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


### Creating PostgreSQL storage

To create and be able to use a PostgreSQL database, a few steps are required:

- In `postgresql.go`, add a record in the list. Commit and push the file; It will spawn the database and create secrets in the `postgres` namespace;
- When database is ready, connect to database with psql to allow user to create tables:

```
$ k exec -ti -n postgres postgres-instance-0 -- psql -U dex-admin dex
dex=# GRANT CREATE ON SCHEMA public TO PUBLIC;
GRANT
```

- Retrieve username & password from `Secret` in the `postgres` namespace:

```sh
$ k get secret -n postgres dex-admin.postgres-instance.credentials.postgresql.acid.zalan.do -o yaml | yq -r .data.username | base64 -d
...
$ k get secret -n postgres dex-admin.postgres-instance.credentials.postgresql.acid.zalan.do -o yaml | yq -r .data.password | base64 -d
...
```

- Create a record in `vault` (ie: `namespaces/<namespace>/postgresql`) and create an `username` and a `password` record with the values retrieved in the `Secret`.

- In the target namespace, create an ExternalSecret

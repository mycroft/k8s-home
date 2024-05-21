# k8s-home

This repository contains source code for my own homelab's kubernetes instance. It is used to define & configure deployed apps, build charts that are then used by gitops mechanisms to deploy them on my cluster.

Thanks to `cdk8s`, the project builds helm charts & apps charts using golang code and stores them in the `generated` branch. This branch is then used by [Flux CD](https://fluxcd.io) to reconcile the k8s state with the built charts.

## Setup

The kubernetes cluster is installed on multiple Mini PCs running Fedora, using [k3s](https://docs.k3s.io/).

## Installed apps

Most important apps installed on my cluster are:

- [longhorn](https://longhorn.io/) is used to provide distributed physical volumes across the cluster;
- [vault](https://www.vaultproject.io/) & [External Secrets Operator](https://external-secrets.io/latest/) are used for most secrets. A few secrets are deployed using [sealed-secrets](https://github.com/bitnami-labs/sealed-secrets), to unlock encrypted volumes;
- [dex-idp](https://dexidp.io/) is linked with my personal [gitea](https://about.gitea.com/) instance to provide oauth based SSO to services supporting it;
- [traefik-forward-auth](https://doc.traefik.io/traefik/middlewares/http/forwardauth/) protects a few services behind an oauth authn/authz firewall for apps not linked to dex-idp;
- [cert-manager](https://cert-manager.io/) generates TLS certificates for ingresses on demand;
- [PostgreSQL](https://www.postgresql.org/), [Minio](https://min.io/), [ScyllaDB](https://www.scylladb.com/), [MariaDB](https://mariadb.org/) operators are installed to provide databases instances;
- [NATS](https://nats.io/) handles message queues;
- [kube-prometheus-stack](https://github.com/prometheus-operator/kube-prometheus), along with [grafana](https://grafana.com/grafana/), provides metrics monitoring;
- Grafana's [loki](https://grafana.com/oss/loki/) with [promtail](https://grafana.com/docs/loki/latest/send-data/promtail/) are setup to aggregate and store logs;
- [linkerd](https://linkerd.io/), [Kyverno](https://kyverno.io/), [Tekton](https://tekton.dev/), [trivy](https://github.com/aquasecurity/trivy) are present in the charts for testing;
- A bunch of end-users apps are installed, such as:
  - [wallabag](https://www.wallabag.it/), a visual to-read list;
  - [freshrss](https://freshrss.org/), a RSS aggregator;
  - [privatebin](https://privatebin.info/), a secure pastebin;
  - [paperless-ngx](https://docs.paperless-ngx.com/) provides document management;
  - [yopass](https://yopass.se/), to share secrets securely;
  - [bookstack](https://www.bookstackapp.com/), a platform to organise and share information;
  - [IT-Tools](https://it-tools.tech/), some handy tools for engineers;
  - [vaultwarden](https://github.com/dani-garcia/vaultwarden), a bitwarden compatible password manager;
  - [send](https://gitlab.com/timvisee/send), simple, private file sharing;
  - [snippetbox](https://github.com/pawelmalak/snippet-box), a code snippet portal;
  - [excalidraw](https://excalidraw.com/), virtual collaborative whiteboard;
  - [wikijs](https://js.wiki/), wiki as stated in its name;
  - [redmine](https://www.redmine.org/), because project management;
  - [microbin](https://microbin.eu/) another pastebin kind of instance.

## Deployment

When a PR is merged in `main`, a Drone CI-based pipeline is started to build cdk8s dependencies, then build charts, and finally push generated charts in the `generated` branch. This branch is pulled every 2 minutes by `flux`. See `.drone.yml`.

## Maintenance notes

### Initial installation with flux

Initial installation of `Flux CD` comes with the `flux` binary. The whole installation procedure is available on the [official website](https://fluxcd.io/flux/get-started/).

Please start `fluxcd`:

```sh
flux bootstrap git \
  --url=ssh://git@git.mkz.me/mycroft/k8s-home.git \
  --private-key-file=/tmp/rsa \
  --branch=generated \
  --path=generated/
```

### Upgrading flux

Dump key file:

```sh
k get secret -n flux-system -o yaml flux-system | yq .data.identity -r | base64 -d > /tmp/rsa
```

Then re-run the `bootstrap` command as seen in the `Installation` section.


### Upgrade services

To check outdated Helm charts & images, use:

```sh
go build && ./k8s-home -check-version
```

Services must be updated by modifying versions in source code. Note all containers are not checked in with this command.

## Security

Secrets can be either stored in source code thanks to `sealed-secrets` or in `vault` and deployed as a `Secret` in Kubernetes with `external-secrets`.

Ingresses can be protected thanks to `traefik-forward-auth` which perform an `Oauth` process using my own private and external `gitea` instance.

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

- In the target namespace, create an ExternalSecretStore and an ExternalSecret. There is a lot of examples in the repository.

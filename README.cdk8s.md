# cdk8s

It seems like `cdk8s` is now overwriting `cdk8s.yaml`, and removes comments.

Imported files are:
- podmonitor:=https://raw.githubusercontent.com/prometheus-community/helm-charts/kube-prometheus-stack-46.3.0/charts/kube-prometheus-stack/crds/crd-podmonitors.yaml
- servicemonitor:=https://raw.githubusercontent.com/prometheus-community/helm-charts/kube-prometheus-stack-46.3.0/charts/kube-prometheus-stack/crds/crd-servicemonitors.yaml
- https://raw.githubusercontent.com/zalando/postgres-operator/v1.9.0/manifests/postgresql.crd.yaml
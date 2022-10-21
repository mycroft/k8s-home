#!/bin/sh

if test "$1" = "extract"; then
  kubectl get secret -n dex-idp dex -o jsonpath='{.data.*}' | base64 -d > current_config.yaml
  cat current_config.yaml
else
  kubectl create secret generic dex-config -n dex-idp --dry-run=client --from-file=foo=current_config.yaml -o json > current_secret.json
  kubeseal < current_secret.json > sealed_secret.json
  cat sealed_secret.json | jq -r ".spec.encryptedData.foo"
fi
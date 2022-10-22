# Vault Notes

## Links

- https://developer.hashicorp.com/vault/docs/platform/k8s/helm
- https://developer.hashicorp.com/vault/docs/platform/k8s/helm/run
- Policies: https://developer.hashicorp.com/vault/docs/concepts/policies
- Policies: https://developer.hashicorp.com/vault/docs/commands/policy

## Setup: Init, unseal

- Install the Helm Chart
- Once installed, init the operator, unseal the vault:

```sh
> kubectl exec vault-0 -- vault operator init -key-shares=1 -key-threshold=1 -format=json > cluster-keys.json
...

> kubectl exec vault-0 -- vault operator unseal <unseal_keys_b64>
...

> kubectl exec vault-0 -- vault status
Key             Value
---             -----
Seal Type       shamir
Initialized     true
Sealed          false
...
```

## Prepare the secret/namespaces/ & create a policy and a token for it:

```sh
> vault secrets enable -path=secret kv-v2
Success! Enabled the kv-v2 secrets engine at: secret/

> vault policy write external-secrets ./external-secrets.vault.hcl
Success! Uploaded policy: external-secrets

```

## Enable kubernetes based authentication

```sh
> vault auth enable kubernetes
Success! Enabled kubernetes auth method at: kubernetes/

> vault write auth/kubernetes/config \
    token_reviewer_jwt="$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
    kubernetes_host="https://$KUBERNETES_PORT_443_TCP_ADDR:443" \
    kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt \
    issuer=https://kubernetes.default.svc
Success! Data written to: auth/kubernetes/config

```


## Allow External-Secrets to access Vault

```sh
> vault write auth/kubernetes/role/external-secrets \
    bound_service_account_names=external-secrets \
    bound_service_account_namespaces=external-secrets \
    policies=external-secrets \
    ttl=60m
Success! Data written to: auth/kubernetes/role/external-secrets

```

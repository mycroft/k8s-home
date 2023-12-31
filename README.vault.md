# Vault Notes

## Links

- https://developer.hashicorp.com/vault/docs/platform/k8s/helm
- https://developer.hashicorp.com/vault/docs/platform/k8s/helm/run
- Auth with Kubernetes: https://developer.hashicorp.com/vault/docs/auth/kubernetes
- kv put: https://developer.hashicorp.com/vault/docs/commands/kv/put
- Policies: https://developer.hashicorp.com/vault/docs/concepts/policies
- Policies: https://developer.hashicorp.com/vault/docs/commands/policy

- External-Secrets: https://external-secrets.io/v0.6.0/provider/hashicorp-vault/
- ES' SecretStore: https://external-secrets.io/v0.6.0/api/secretstore/

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

> vault policy write external-secrets-ui ./external-secrets-ui.vault.hcl
Success! Uploaded policy: external-secrets-ui

```

## Enable kubernetes based authentication

In the past, the serviceaccount/token was passed in `token_reviewer_jwt` and stored in vault. However, this breaks Vault on restart.
The `auth/kubernetes/config` should only contain the `kubernetes_host` API URL. It was tested as working after a container restart.

```sh
> vault auth enable kubernetes
Success! Enabled kubernetes auth method at: kubernetes/

> vault write auth/kubernetes/config kubernetes_host=https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT
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

## Enable userpass

```sh
> vault auth enable userpass
Success! Enabled userpass auth method at: userpass/

> vault write auth/userpass/users/mycroft policies=default,external-secrets-ui password=IghoPoh9eech/aca
Success! Data written to: auth/userpass/users/mycroft

```

## Write secrets

```sh
> vault kv put secret/namespaces/external-secrets/testaroo username="static-user" password="static-password" field="new-field"

```
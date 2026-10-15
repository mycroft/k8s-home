# Vault Notes

## Links

- https://developer.hashicorp.com/vault/docs/platform/k8s/helm
- https://developer.hashicorp.com/vault/docs/platform/k8s/helm/run
- Auth with Kubernetes: https://developer.hashicorp.com/vault/docs/auth/kubernetes
- kv put: https://developer.hashicorp.com/vault/docs/commands/kv/put
- Policies: https://developer.hashicorp.com/vault/docs/concepts/policies
- Policies: https://developer.hashicorp.com/vault/docs/commands/policy
- OIDC auth: https://developer.hashicorp.com/vault/docs/auth/jwt
- Authentik's Vault guide: https://integrations.goauthentik.io/security/hashicorp-vault/

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

## Log in to the container

Unlike the init and unseal steps above, every `vault` command that follows runs
from inside the Vault pod and needs an authenticated session:

```sh
> kubectl exec -ti -n vault vault-0 -- /bin/sh
/ $ vault login
Token (will be hidden):
```

The root token is stored in the [pass(1)](https://www.passwordstore.org/)
password manager.

The policy files live in `contrib/` in this repository and are not present in
the container, so copy them over before applying them:

```sh
> kubectl cp contrib/external-secrets.vault.hcl vault/vault-0:/tmp/
> kubectl cp contrib/external-secrets-ui.vault.hcl vault/vault-0:/tmp/
```

## Prepare the secret/namespaces/ & create a policy and a token for it:

```sh
> vault secrets enable -path=secret kv-v2
Success! Enabled the kv-v2 secrets engine at: secret/

> vault policy write external-secrets /tmp/external-secrets.vault.hcl
Success! Uploaded policy: external-secrets

> vault policy write external-secrets-ui /tmp/external-secrets-ui.vault.hcl
Success! Uploaded policy: external-secrets-ui

```

Re-run `vault policy write` the same way whenever a policy in `contrib/`
changes: the files are not managed by Flux, so nothing applies them for you.
Because tokens refer to policies by name, the updated rules apply to existing
tokens immediately.

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

Day-to-day login goes through Authentik (next section); `userpass` is the
break-glass path and must stay enabled. Authentik reads its own secrets from
Vault (ExternalSecrets in `charts/security/authentik.go`), so when Vault is
sealed or Authentik is down, OIDC login is unavailable exactly when it is needed.

```sh
> vault auth enable userpass
Success! Enabled userpass auth method at: userpass/

> vault write auth/userpass/users/mycroft policies=default,external-secrets-ui password=some-password
Success! Data written to: auth/userpass/users/mycroft

```

## Log in with Authentik (OIDC)

### Authentik

Created by hand in the admin UI, like the other apps (see `authentik.md`):

- Group `Vault Admins`, holding the users who get `external-secrets-ui`.
- OAuth2/OpenID provider `vault`:
  - Client type: Confidential.
  - Redirect URIs, both Strict and of type Authorization:
    - `https://vault.services.mkz.me/ui/vault/auth/oidc/oidc/callback` (UI)
    - `http://localhost:8250/oidc/callback` (`vault login -method=oidc`)
  - Signing Key: pick a certificate. Without one, authentik signs ID tokens
    with HS256 and serves an empty JWKS, which Vault cannot verify.
  - Scopes: keep the defaults (`openid`, `email`, `profile`). The `groups`
    claim comes from `profile`.
- Application `Vault`, slug `vault`, using that provider.

Check the signing key from outside: this should print `["RS256"]`.

```sh
> curl -s https://auth.services.mkz.me/application/o/vault/.well-known/openid-configuration | jq .id_token_signing_alg_values_supported
```

### Vault

From the container, logged in with the root token:

```sh
/ $ vault auth enable oidc
Success! Enabled oidc auth method at: oidc/

/ $ vault auth tune -listing-visibility=unauth oidc/
Success! Tuned the auth method at: oidc/
```

`-listing-visibility=unauth` puts the OIDC method on the UI login page.

The client secret goes in on stdin (`=-`) to keep it out of the shell history.
Paste it, press Enter, then Ctrl-D. `tr` drops the newline, which Vault would
otherwise keep as part of the secret:

```sh
/ $ tr -d '\n' | vault write auth/oidc/config \
      oidc_discovery_url="https://auth.services.mkz.me/application/o/vault/" \
      oidc_client_id="<client id>" \
      oidc_client_secret=- \
      default_role=default
Success! Data written to: auth/oidc/config

/ $ vault write auth/oidc/role/default \
      role_type=oidc \
      bound_audiences="<client id>" \
      allowed_redirect_uris="https://vault.services.mkz.me/ui/vault/auth/oidc/oidc/callback" \
      allowed_redirect_uris="http://localhost:8250/oidc/callback" \
      user_claim=preferred_username \
      groups_claim=groups \
      oidc_scopes="profile,email" \
      token_policies=default
Success! Data written to: auth/oidc/role/default
```

Vault always requests `openid`, so `oidc_scopes` lists only the extra scopes.
The role grants nothing beyond `default`; other policies come from groups.

Map the Authentik group to a Vault external group. The alias name must match
the group name in the `groups` claim exactly. Writing the group by name is safe
to repeat; the container has no `jq`, so ids are read with `-field`:

```sh
/ $ vault write identity/group/name/"Vault Admins" type=external policies=external-secrets-ui
/ $ GID=$(vault read -field=id identity/group/name/"Vault Admins")
/ $ ACCESSOR=$(vault read -field=accessor sys/auth/oidc)
/ $ vault write identity/group-alias name="Vault Admins" mount_accessor="$ACCESSOR" canonical_id="$GID"
```

Vault refreshes external group membership on login and token renewal, so
removing someone from `Vault Admins` takes effect then. Revoke their token to
cut access sooner.

### Logging in

- UI: choose the OIDC method on https://vault.services.mkz.me and leave Role
  empty (`default_role` applies).
- CLI: run it from a workstation, not the container. The CLI listens on
  `localhost:8250` for the callback and opens a browser.

```sh
> export VAULT_ADDR=https://vault.services.mkz.me
> vault login -method=oidc
> vault token lookup
```

A `Vault Admins` member gets `policies [default]` and
`identity_policies [external-secrets-ui]`. Anyone else gets only `default`.

## Write secrets

```sh
> vault kv put secret/namespaces/external-secrets/testaroo username="static-user" password="static-password" field="new-field"

```

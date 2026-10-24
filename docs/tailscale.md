# Tailscale

The [Tailscale Kubernetes operator](https://tailscale.com/docs/kubernetes-operator)
(`charts/infra/tailscale.go`) exposes selected services on the tailnet only. For every
Ingress of class `tailscale` it runs a proxy pod in the `tailscale` namespace that joins
the tailnet as its own machine, `<host>.<tailnet>.ts.net`, and serves HTTPS there with a
Let's Encrypt certificate it obtains itself. Traefik and cert-manager are not involved, and
the host cannot be reached from the internet.

In a chart:

```go
chart.NewTailscaleIngress("vaultwarden", svcName, appPort)
// → https://vaultwarden.<tailnet>.ts.net
```

## One-time tailnet setup

All of it happens in the admin console, as an Owner or Admin of the tailnet. Do it before
the chart is deployed: the operator pod cannot start until its credentials exist.

1. **Tags** — *Access controls* (<https://console.tailscale.com/admin/acls>). In the JSON
   editor, add to `tagOwners` (create the section if it is missing):

   ```json
   "tagOwners": {
     "tag:k8s-operator": ["autogroup:admin"],
     "tag:k8s":          ["tag:k8s-operator"],
   },
   ```

   The operator authenticates as `tag:k8s-operator` and gives every proxy `tag:k8s`; both
   are pinned in `configs/tailscale-operator.yaml`. If the policy file has its own
   `grants` or `acls` (not the default allow-all), also let your devices reach the proxies:

   ```json
   "grants": [
     {"src": ["autogroup:member"], "dst": ["tag:k8s"], "ip": ["443"]},
   ],
   ```

2. **MagicDNS and HTTPS** — *DNS* (<https://console.tailscale.com/admin/dns>). Enable
   MagicDNS, then under *HTTPS Certificates* select *Enable HTTPS*. The tailnet name shown
   on that page is the `<tailnet>` part of every `*.ts.net` hostname — here,
   `merino-tiaki.ts.net`.

3. **OAuth client** — *Settings › Trust credentials*
   (<https://console.tailscale.com/admin/settings/trust-credentials>; older docs call these
   "OAuth clients"). Select *+ Credential*, then *OAuth*, and grant Read and Write on:
   - General › Services
   - Devices › Core
   - Keys › Auth Keys

   Under *Tags*, add `tag:k8s-operator` (it is only selectable once step 1 is saved).
   Select *Generate credential* and copy the client ID and secret: the secret is shown
   once.

4. **Vault** — store the credentials where the chart's ExternalSecret reads them:

   ```sh
   vault kv put secret/namespaces/tailscale/operator-oauth \
       client_id=<client id> client_secret=<client secret>
   ```

   They become the `operator-oauth` Secret, which the Helm chart mounts when `oauth.*` is
   left unset.

## Checking it works

```sh
kubectl -n tailscale get pods                    # operator, plus one ts-* proxy per Ingress
kubectl -n vaultwarden get ingress vaultwarden-tailscale   # ADDRESS: vaultwarden.merino-tiaki.ts.net
```

The operator appears in the admin console's *Machines* page as `tailscale-operator`, each
proxy under its Ingress host. The first HTTPS request to a new host can take up to a
minute while its certificate is issued.

## Troubleshooting

- **Operator pod stuck in `ContainerCreating`** — the `operator-oauth` Secret is missing.
  Check the Vault entry and `kubectl -n tailscale get externalsecret`; the ExternalSecret
  refreshes every 15 minutes.
- **`requested tags [...] are invalid or not permitted`** in
  `kubectl -n tailscale logs deploy/operator` — the `tagOwners` from step 1 are missing,
  or the OAuth client was not given `tag:k8s-operator`.
- **Host resolves but the connection is refused or times out** — a custom ACL is
  blocking `tag:k8s` (step 1), or HTTPS is not enabled for the tailnet (step 2).
- **Proxy pods run a privileged init container** that enables IP forwarding. The operator
  creates them at runtime, so they never appear in `dist/` and conftest does not see them;
  Falco may report them.

## Vaultwarden migration

Vaultwarden is currently served both publicly (`vaultwarden.services.mkz.me`, through
Traefik) and on the tailnet. `DOMAIN` still names the public host. To finish the move:

1. Point each client at `https://vaultwarden.merino-tiaki.ts.net` (browser extension, mobile
   and desktop apps: *self-hosted* server URL on the login screen). Changing the server
   logs the client out. While `DOMAIN` is the public host, security-key (FIDO2/WebAuthn)
   two-factor and passkey login only work there: log in with another second factor, or
   switch to TOTP first.
2. Once every client has moved, in `charts/apps/vaultwarden.go` set `DOMAIN` to
   `https://vaultwarden.merino-tiaki.ts.net` and remove the `NewAppIngress` call.
3. Register security keys and passkeys again: they are bound to the domain. The vault data
   itself is not.

Clients keep a cached copy of the vault and can read it offline, but sync and edits need
the tailnet. If the Tailscale login (or the identity provider behind it) lives in
Vaultwarden, keep a client logged in or an offline export, so a tailnet outage cannot lock
you out of both.

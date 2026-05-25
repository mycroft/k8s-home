# Notes about authentik

## Setup

Authentik is not using default PostgreSQL user but the admin's as it must create schemas and so on.

## Application setup

There is a default outpost and middleware for it created in this repository.

### Complete app setup:

Setting up Authentik:

- Create an Outpost Integration
  - Name: Local Kubernetes Outpost Integration
  - Local

- Create a new Provider:
  - Name: Forward Auth for Bookstack
  - Authn flow: default-authentification-flow
  - Authz flow: default-provider-authorization-explicit-consent
  - Forward auth (single app)
  - External host: https://bookstack.services.mkz.me/

- Create an application:
  - Name: Bookstack
  - slug: bookstack
  - Provider: Forward Auth for Bookstack
  - Policy engine mode: any

- Create an outpost:
  - Name: Outpost for Bookstack
  - Type: Proxy
  - Integration: Local Kubernetes Outpost Integration
  - Apps: Bookstack

- Create an IngressRoute matching the newly created service:

```
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: bookstack
spec:
  routes:
  - kind: Rule
    match: "Host(`bookstack.services.mkz.me`)"
    middlewares:
    - name: ak-outpost-outpost-for-bookstack
      namespace: authentik
    priority: 10
    services:
    - name: bookstack-svc-c8e87854
      port: http
  tls:
    secretName: secret-tls-www
```

If you want to use default outpost, just use `authentik` middleware in `authentik` namespace.

## Misc

To enforce MFA:
- Edit the default `default-authentication-mfa-validation stage and set the Not configured action to Force the user to configure an authenticator and select a configuration stage for the device type you want (for example default-authenticator-totp-setup)
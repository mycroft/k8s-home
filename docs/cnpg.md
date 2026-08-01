# CloudNativePG

A single CloudNativePG cluster named `postgres` runs in the `cnpg` namespace and
hosts one database per application. It is declared in
`charts/storage/cnpg-cluster.go`; the per-database plumbing lives in
`internal/kubehelpers/cnpg.go`.

Applications reach it at `postgres-rw.cnpg`, the read-write service.

## Adding a new database

Add one entry to the `databases` list in `charts/storage/cnpg-cluster.go`:

```go
var databases = []kubehelpers.CNPGDatabase{
    {Name: "zipline", VaultEntry: "postgresql-cnpg"},
    {Name: "myapp"},
}
```

`Name` is used for the database, its owning role and the generated secrets, and
is the only required field:

| Field        | Default          | Description                                              |
| ------------ | ---------------- | -------------------------------------------------------- |
| `Name`       | —                | Database name, owning role, and secret name prefix       |
| `Namespace`  | `Name`           | Application namespace the credentials are published to   |
| `VaultEntry` | `postgresql`     | Secret name under that namespace's Vault path            |

Set `Namespace` when the application namespace differs from the database name.
Set `VaultEntry` to stage credentials on a separate Vault path while migrating,
so a still-running application's existing credentials are not overwritten.

Regenerate and commit as usual — the role is part of the cluster spec, so the
same entry both adds the managed role and creates the database resources:

```sh
mise run generate-charts
```

## Where the secrets are created

Everything below is named `postgres-<name>` and created in the `cnpg` namespace:

| Resource                     | Purpose                                                     |
| ---------------------------- | ----------------------------------------------------------- |
| `Password` generator         | Mints a random 32-character password, no symbols            |
| `ExternalSecret`             | Materializes it into a `kubernetes.io/basic-auth` secret    |
| `Secret`                     | Consumed by the cluster's managed role as `passwordSecret`  |
| `Database`                   | Creates the database on the cluster, owned by the role      |
| `PushSecret`                 | Publishes the credentials to Vault                          |

The secret holds five keys, and the same five are pushed to Vault:

| Key        | Value                                            |
| ---------- | ------------------------------------------------ |
| `username` | the database name                                |
| `password` | the generated password                           |
| `host`     | `postgres-rw.cnpg`                               |
| `dbname`   | the database name                                |
| `url`      | `postgres://<name>:<password>@postgres-rw.cnpg/<name>` |

The Vault path is `namespaces/<namespace>/<vault entry>` relative to the KV
mount, so `{Name: "myapp"}` publishes to `secret/namespaces/myapp/postgresql` —
exactly where `CreateExternalSecret` looks. Nothing else is needed to hand the
credentials over to the application:

```go
kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql")
```

That creates a `postgresql` secret in the application namespace, whose `url` key
is usually wired straight into the container:

```go
{
    Name: jsii.String("DATABASE_URL"),
    ValueFrom: &k8s.EnvVarSource{
        SecretKeyRef: &k8s.SecretKeySelector{
            Key:  jsii.String("url"),
            Name: jsii.String("postgresql"),
        },
    },
},
```

## Notes

The password is pinned with `refreshInterval: 0`: without it every refresh would
mint a new one and lock the application out until it restarts. Editing the
`ExternalSecret` template in `cnpgCredentialsTemplate` re-reconciles it and does
regenerate the password, so treat changes there as a credential rotation.

Passwords are generated without symbols so they stay safe to embed in the `url`
without escaping.

The `PushSecret` needs both `secret/data/namespaces/*` and
`secret/metadata/namespaces/*` in the `external-secrets` Vault policy — it
maintains the KV-v2 metadata of what it publishes alongside the data. See
[docs/vault.md](vault.md) for applying policy changes.

Removing an entry from the list deletes its Kubernetes resources, including the
generated credential Secret. The database itself survives because
`databaseReclaimPolicy` defaults to `retain`, and the Vault copy survives because
the PushSecret uses `deletionPolicy: None`. The retained database and Vault
credentials have to be dropped by hand.

# CloudNativePG

A single CloudNativePG cluster named `postgres` runs in the `cnpg` namespace and
hosts one database per application. It is declared in
`charts/storage/cnpg-cluster.go`; the per-database plumbing lives in
`internal/kubehelpers/cnpg.go`.

Applications reach it at `postgres-rw.cnpg`, the read-write service.

This replaces the former Zalando PostgreSQL operator. There are no
`acid.zalan.do` resources or Zalando operator namespaces in the current
deployment.

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

| Field             | Default      | Description                                                   |
| ----------------- | ------------ | ------------------------------------------------------------- |
| `Name`            | —            | Database name, owning role, and secret name prefix            |
| `Namespace`       | `Name`       | Application namespace the credentials are published to        |
| `VaultEntry`      | `postgresql` | Secret name under that namespace's Vault path                 |
| `PasswordAliases` | none         | Additional keys containing the generated password             |
| `Extensions`      | none         | PostgreSQL extensions that CNPG must keep installed            |

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

The secret holds five core keys, and the same keys are pushed to Vault:

| Key        | Value                                            |
| ---------- | ------------------------------------------------ |
| `username` | the database name                                |
| `password` | the generated password                           |
| `host`     | `postgres-rw.cnpg`                               |
| `dbname`   | the database name                                |
| `url`      | `postgres://<name>:<password>@postgres-rw.cnpg/<name>` |

Every configured `PasswordAliases` entry is added to the Secret and Vault with
the same value as `password`. This supports applications whose charts require a
specific password key without duplicating credentials.

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

## Backups

Database recovery uses logical PostgreSQL dumps backed by Velero rather than a
raw copy of a running PostgreSQL data directory:

1. The `postgres-dump` CronJob starts daily at `00:30` UTC in `cnpg`.
2. It writes `globals.sql` plus one custom-format `<database>.dump` file per
   connectable database into a timestamped directory on the `postgres-dumps`
   PVC.
3. A completed directory is atomically exposed through `/backup/latest`.
4. The CronJob keeps the PVC mounted while the Velero filesystem schedule runs
   at `01:30` UTC.
5. A Velero pre-backup hook rejects empty, incomplete, or stale dumps while the
   dump pod is active.
6. Kopia copies the dump PVC and every other mounted persistent volume to the
   configured Velero object store.

Both local logical dumps and Velero backups are retained for 30 days. Raw CNPG
volumes are included by the cluster-wide filesystem schedule, but logical dumps
are the supported database restore source because copying a live PostgreSQL data
directory does not provide the same consistency guarantee.

### Verify the daily backup

Check the CronJob and its latest Job:

```sh
kubectl get cronjob postgres-dump -n cnpg
kubectl get jobs -n cnpg -l app.kubernetes.io/name=postgres-dump \
    --sort-by=.metadata.creationTimestamp
```

While the latest dump pod is running during its eight-hour backup window,
verify the directory and non-empty dump files:

```sh
kubectl exec -n cnpg job/<job-name> -- sh -c \
    'readlink /backup/latest && find -L /backup/latest -maxdepth 1 -type f -size +0c -printf "%f %s bytes\n" | sort'
```

Then check the corresponding Velero backup:

```sh
kubectl get backups.velero.io -n velero \
    --sort-by=.metadata.creationTimestamp
kubectl describe backup.velero.io -n velero <backup-name>
```

The filesystem backup must be `Completed`, and the dump Job must exist and have
produced `/backup/latest`; without a dump pod there is no hook or mounted dump
PVC for Velero to inspect. A failed `postgres-dump-ready` hook means Velero did
not find a complete dump produced within the previous three hours. Fix the
CronJob before accepting that backup.

### Run a dump manually

Create a one-off Job from the deployed CronJob:

```sh
kubectl create job -n cnpg --from=cronjob/postgres-dump \
    postgres-dump-manual-$(date -u +%Y%m%d%H%M%S)
```

The dump itself normally finishes quickly, but the Job deliberately remains
running for eight hours so Kopia can access the mounted PVC. Inspect the files
using the command above; do not wait for Job completion to confirm the dump.

### Restore a database

Restore the `postgres-dumps` PVC from Velero into an isolated recovery
environment first. Stop writers before restoring into an existing database.
Apply global roles before restoring database contents:

```sh
psql --host=<recovery-host> --username=postgres --dbname=postgres \
    --file=/backup/<timestamp>/globals.sql

pg_restore --host=<recovery-host> --username=postgres \
    --dbname=<database> --clean --if-exists \
    /backup/<timestamp>/<database>.dump
```

The target database must already exist. In this repository the CNPG `Database`
resource normally creates it with the correct owner. Validate the recovered
application before allowing traffic or removing the pre-restore backup.

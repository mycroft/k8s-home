# Garage Setup and Maintenance Playbook

## Executive Summary

Garage runs as a single-node S3-compatible object store on `moonstone`. It currently provides a `300G` layout in zone `dc1` and a `backup` bucket, but the single-node deployment has no storage redundancy and must not be treated as its own backup.

This playbook records the known setup and the procedures used to verify, maintain, upgrade, and recover it. Run Garage administration commands on `moonstone` as a user that can read the Garage configuration and metadata directory.

## Current Deployment

| Item | Value |
| --- | --- |
| Host | `moonstone.lan.mkz.me` |
| Node ID | `1b2403d9bf411542` |
| Zone | `dc1` |
| Allocated capacity | `300G` |
| Web endpoint | `http://moonstone.lan.mkz.me:3902/garage` |
| Admin/metrics endpoint | `moonstone.lan.mkz.me:3903` |
| Backup bucket | `backup` |
| Backup key name | `backup-key` |
| TLS certificate | `s3.lan.mkz.me`, synced from the cluster to `/etc/garage-tls/` (see [TLS Certificate](#tls-certificate)) |

Garage conventionally uses port `3900` for its S3 API, `3901` for RPC, `3902` for web access, and `3903` for its admin API. Confirm the actual bindings and data paths in the active Garage configuration before rebuilding the service:

```sh
garage --version
garage status
garage layout show
```

The default configuration path is `/etc/garage.toml`. If a different path is used, export it before running commands:

```sh
export GARAGE_CONFIG_FILE=/path/to/garage.toml
```

For a container deployment, run the same commands inside the Garage container.

## Initial Layout

The observed initial setup was:

```sh
garage layout assign -z dc1 1b2403d9bf411542
garage layout assign -z dc1 -c 300G 1b2403d9bf411542
garage layout apply --version 1
garage status
```

For a clean rebuild, the two assignment commands can be reduced to one:

```sh
garage status
garage layout assign -z dc1 -c 300G 1b2403d9bf411542
garage layout show
garage layout apply --version 1
garage status
```

Only use `--version 1` for a new layout. On an existing cluster, inspect `garage layout show` and apply the staged layout version shown by Garage. Applying the wrong version is rejected and should not be worked around by guessing.

After applying the layout, `garage status` should report:

- Node `1b2403d9bf411542` as healthy
- Zone `dc1`
- Capacity `300G`
- No nodes under `FAILED NODES`

## Backup Bucket and Credentials

The `backup` bucket and its access key were created with:

```sh
garage bucket create backup
garage key create backup-key
garage bucket allow --read --write --owner --key backup-key backup
```

The key creation command prints the secret key. Store it immediately in Vault or another password manager; do not add it to this repository.

Verify the authorization:

```sh
garage bucket info backup
garage key info backup-key
```

The bucket should list `backup-key` with read, write, and owner permissions.

### S3 Client Test

Assuming the standard S3 binding on port `3900`, configure a temporary shell without writing credentials to disk:

```sh
export AWS_ENDPOINT_URL=http://moonstone.lan.mkz.me:3900
export AWS_DEFAULT_REGION=garage
export AWS_ACCESS_KEY_ID=<backup-key-id>
export AWS_SECRET_ACCESS_KEY=<backup-secret-key>
```

Exercise both write and read paths:

```sh
test_file=$(mktemp)
printf 'garage smoke test\n' > "$test_file"
aws s3 cp "$test_file" s3://backup/maintenance/smoke-test.txt
aws s3 cp s3://backup/maintenance/smoke-test.txt -
aws s3 rm s3://backup/maintenance/smoke-test.txt
rm "$test_file"
```

The upload, download, and removal must all succeed. Unset credentials afterward:

```sh
unset AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_ENDPOINT_URL AWS_DEFAULT_REGION
```

## TLS Certificate

Garage does not serve TLS on its S3 API, so a proxy on moonstone terminates it. The certificate comes from the cluster:

1. cert-manager issues it from `letsencrypt-prod` into Secret `s3-tls` in namespace `moonstone-garage` (`charts/storage/moonstone-garage.go`). `lan.mkz.me` names resolve to private addresses Let's Encrypt cannot reach, so both issuers validate them with DNS-01: cert-manager writes the `_acme-challenge` TXT record to maki's BIND with TSIG key `cert-manager` (Vault: `secret/namespaces/cert-manager/rfc2136-tsig`, field `secret`), which `update-policy` restricts to those names.
2. On moonstone, `garage-tls-sync.timer` runs `garage-tls-sync` hourly (`contrib/moonstone-garage/`). It reads the Secret as ServiceAccount `moonstone`, which may `get` that one Secret and nothing else, installs `tls.crt` (full chain) and `tls.key` into `/etc/garage-tls/`, then runs `RELOAD_COMMAND`. Files are only replaced when the Secret changed, and an expired certificate or a mismatched key is refused.
3. The proxy serves `/etc/garage-tls/tls.crt` and `/etc/garage-tls/tls.key`. Until it is in place, the S3 API stays plain HTTP on port `3900`.

cert-manager renews when a third of the certificate's lifetime remains; `kubectl -n moonstone-garage get certificate s3-tls -o jsonpath='{.status.renewalTime}'` shows when. If the cluster is down, moonstone keeps serving the installed certificate until it expires.

### Setting Up the Sync

Moonstone needs `curl`, `jq`, `openssl`, and access to `star-train.lan.mkz.me:6443`.

From a workstation with cluster access, copy the ServiceAccount token, the cluster CA and the files from `contrib/moonstone-garage/`:

```sh
kubectl -n moonstone-garage get secret moonstone-token -o jsonpath='{.data.token}' | base64 -d \
  | ssh moonstone 'umask 077; cat > garage-tls-token'
kubectl -n moonstone-garage get secret moonstone-token -o jsonpath='{.data.ca\.crt}' | base64 -d \
  | ssh moonstone 'cat > garage-tls-ca.crt'
scp contrib/moonstone-garage/garage-tls-sync* moonstone:
```

Then on moonstone:

```sh
sudo install -d -m 0755 /etc/garage-tls
sudo install -d -m 0700 /etc/garage-tls/sa
sudo install -m 0600 garage-tls-token /etc/garage-tls/sa/token
sudo install -m 0644 garage-tls-ca.crt /etc/garage-tls/sa/ca.crt
sudo install -m 0755 garage-tls-sync /usr/local/bin/
sudo install -m 0644 garage-tls-sync.service garage-tls-sync.timer /etc/systemd/system/
sudo install -m 0644 garage-tls-sync.env /etc/default/garage-tls-sync
rm garage-tls-token garage-tls-ca.crt garage-tls-sync*

sudo systemctl daemon-reload
sudo systemctl start garage-tls-sync.service
journalctl -u garage-tls-sync -n 5
sudo systemctl enable --now garage-tls-sync.timer
```

The first run logs `s3-tls: installed certificate, notAfter=...`. Get to that point before pointing the proxy at the files: `nginx -t` fails while they do not exist. Then set `RELOAD_COMMAND` in `/etc/default/garage-tls-sync` (and `KEY_GROUP` if the proxy does not run as root).

### Verifying

```sh
# on moonstone
sudo openssl x509 -in /etc/garage-tls/tls.crt -noout -subject -issuer -enddate
systemctl list-timers garage-tls-sync.timer

# from the workstation
kubectl -n moonstone-garage get certificate s3-tls
```

To exercise a renewal end to end, run `cmctl renew -n moonstone-garage s3-tls`, wait for the Certificate to be Ready again, then `sudo systemctl start garage-tls-sync.service` on moonstone and check that `notAfter` moved. Forced renewals count against Let's Encrypt's rate limits; do not loop on it.

### Revoking Moonstone's Access

```sh
kubectl -n moonstone-garage delete secret moonstone-token
```

The old token stops working immediately. Flux recreates the Secret with a new token on its next reconcile; install that on moonstone as in [Setting Up the Sync](#setting-up-the-sync).

## Routine Maintenance

### Weekly

Check node health, layout, disk usage, buckets, and keys:

```sh
garage status
garage layout show
garage bucket list
garage key list
```

Investigate unexpected failed nodes, layout changes, or a rapid reduction in available capacity before accepting further backup traffic.

The Prometheus target for Garage is configured at `moonstone.lan.mkz.me:3903`. Confirm it is up in Prometheus and review Garage capacity and error metrics.

### Monthly

Perform the S3 client test above and verify that at least one real backup consumer can:

1. Write a new backup
2. List it through the S3 API
3. Restore or download it
4. Validate the restored content

Listing objects alone does not prove that backup data is complete or recoverable.

### Changing Capacity

First inspect the active and staged layout:

```sh
garage layout show
```

Stage the new capacity:

```sh
garage layout assign -z dc1 -c <new-capacity> 1b2403d9bf411542
garage layout show
```

Apply the version displayed for the staged layout:

```sh
garage layout apply --version <staged-version>
garage status
```

Do not advertise more capacity than the filesystem can safely provide. Leave enough free space for compaction, temporary files, operating-system updates, and recovery work.

### Rotating an Application Key

Create a replacement before removing the old key:

```sh
garage key create <new-key-name>
garage bucket allow --read --write --owner --key <new-key-name> backup
garage key info <new-key-name>
```

Store the new secret, update all clients, and verify a complete write/read/delete test. Only then revoke or delete the old key using the command shown by:

```sh
garage key --help
```

## Upgrades

Garage is a single-node service, so upgrades cause an availability interruption.

Before upgrading:

1. Read the release notes for every skipped version.
2. Record `garage --version`, `garage status`, and `garage layout show`.
3. Back up `/etc/garage.toml`, the RPC secret, administration tokens, and all access credentials.
4. Take a cold backup of the configured metadata and data directories, or verify a current independent object-level replica.
5. Stop all backup writers.

Upgrade the binary or container using the host's existing service-management method, then start Garage and verify:

```sh
garage --version
garage status
garage layout show
garage bucket info backup
garage key info backup-key
```

Finish with the S3 client test. Do not remove the previous binary, container image, configuration backup, or data backup until the service has passed these checks.

## Backup and Recovery

The Moonstone deployment is a single failure domain. Regardless of the configured replication factor, a single physical node does not protect against disk loss or loss of the host.

Maintain an independent copy using one or both of these approaches:

- Replicate the `backup` bucket to an external S3 provider using `rclone`, with encryption enabled on the destination.
- Stop Garage and take a consistent host-level backup of its configuration, metadata directory, and data directory.

Keep copies of the following outside Moonstone:

- Garage configuration
- RPC secret
- Admin and metrics tokens
- Access key IDs and secret keys
- Metadata directory
- Object data or an independently replicated bucket

### Recovery Checklist

1. Provision a replacement host and storage.
2. Restore the exact Garage version used by the backup when possible.
3. Restore the configuration, metadata, and data directories with their original ownership and permissions.
4. Start Garage.
5. Run `garage status` and `garage layout show`.
6. Verify `backup` and `backup-key`.
7. Perform the S3 client test.
8. Restore an actual backup consumer into an isolated environment.

If only an object-level replica is available, create a fresh Garage instance and copy the objects back through the S3 API rather than attempting to reconstruct Garage's internal metadata manually.

## Repository Follow-up

The repository still registers `charts/storage/garage.go`, which describes the former in-cluster Garage deployment. Confirm whether that release is intentionally still active before treating Moonstone as the only Garage instance. Removing or disabling the Kubernetes release should be handled as a separate, explicitly reviewed migration.

## Troubleshooting

### CLI Cannot Connect

- Confirm Garage is running.
- Confirm the CLI reads the active configuration file.
- Check the RPC address and RPC secret.
- Run the CLI locally on `moonstone` before testing remotely.

### Node Has No Assigned Role

```sh
garage status
garage layout show
garage layout assign -z dc1 -c 300G <node-id>
garage layout show
garage layout apply --version <staged-version>
```

### Certificate Sync Fails

Start with `journalctl -u garage-tls-sync` on moonstone. A failed run never touches the installed certificate.

- **`/etc/garage-tls/sa/token is missing`**: the token and CA were never installed; see [Setting Up the Sync](#setting-up-the-sync).
- **`curl: (22) ... 401`**: the token is no longer valid, usually because `moonstone-token` was recreated. Install the new one.
- **`curl: (22) ... 403`**: the `s3-tls-reader` Role or RoleBinding no longer grants access.
- **`curl: (7)` or a timeout**: the API server at `star-train.lan.mkz.me:6443` is unreachable.
- **`Secret has no tls.crt/tls.key yet`**: the certificate was never issued. Check `kubectl -n moonstone-garage describe certificate s3-tls`, then its `challenge`. If the challenge reports a DNS update error, named's log on maki shows whether the update was denied (key name, algorithm or `update-policy` grant mismatch).

### S3 Requests Fail

- Confirm the client uses the S3 API endpoint, normally port `3900`, rather than the web endpoint on `3902`.
- Confirm the region matches the configured Garage S3 region.
- Confirm `backup-key` is authorized for the `backup` bucket.
- Confirm the client uses path-style addressing when required.
- Check Garage logs with `RUST_LOG=garage=debug` only while diagnosing; return to `info` afterward.

## References

- [Garage documentation](https://garagehq.deuxfleurs.fr/documentation/)
- [Garage configuration reference](https://garagehq.deuxfleurs.fr/documentation/reference-manual/configuration/)
- [Garage administration commands](https://garagehq.deuxfleurs.fr/documentation/reference-manual/administration/)

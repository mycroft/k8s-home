# This must be used to create tokens to be used for external-secrets
#
# # vault policy write external-secrets ./external-secrets.vault.hcl
# Success! Uploaded policy: external-secrets

path "secret/data/namespaces/*" {
  capabilities = ["create", "read", "update", "patch", "delete", "list"]
}

# PushSecret maintains the KV-v2 metadata of the secrets it publishes alongside
# their data, so writing credentials back to Vault needs this path too.
path "secret/metadata/namespaces/*" {
  capabilities = ["create", "read", "update", "patch", "delete", "list"]
}

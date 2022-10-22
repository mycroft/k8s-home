# This must be used to create tokens to be used for external-secrets
#
# # vault policy write external-secrets ./external-secrets.vault.hcl
# Success! Uploaded policy: external-secrets

path "secret/namespaces/*" {
  capabilities = ["create", "read", "update", "patch", "delete", "list"]
}

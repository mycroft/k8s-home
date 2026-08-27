package main

import rego.v1

# vault_keys is the set of Vault paths an ExternalSecret references. Only
# extract entries carry paths; generator-sourced entries (e.g. the CNPG
# Password generators) have no Vault path and are out of scope.
vault_keys contains k if {
	k := input.spec.dataFrom[_].extract.key
}

vault_keys contains k if {
	k := input.spec.data[_].remoteRef.path
}

# Rule ID: vault-path
# ExternalSecrets must extract from the fleet's Vault layout:
# secret/namespaces/<this object's namespace>/<secret-name>.
# The standard helper (CreateExternalSecret) builds exactly this path, so a
# violation means a hand-written ExternalSecret or a mismatched namespace
# argument — e.g. a copy-pasted ES extracting from another namespace's path.
deny contains msg if {
	input.kind == "ExternalSecret"
	some k in vault_keys
	not startswith(k, sprintf("secret/namespaces/%s/", [input.metadata.namespace]))
	msg := sprintf("ExternalSecret '%v' in '%v' extracts from '%v'; expected a path under secret/namespaces/%v/", [input.metadata.name, input.metadata.namespace, k, input.metadata.namespace])
}

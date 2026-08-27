package main

import rego.v1

# storage_claims is the set of {name, class} pairs for the storage an
# object requests: the PVC itself, or each StatefulSet volumeClaimTemplate.
# A missing storageClassName is reported as <none>.
storage_claims contains e if {
	input.kind == "PersistentVolumeClaim"
	e := {"name": input.metadata.name, "class": object.get(input.spec, "storageClassName", "<none>")}
}

storage_claims contains e if {
	input.kind == "StatefulSet"
	some v in input.spec.volumeClaimTemplates
	e := {"name": v.metadata.name, "class": object.get(v.spec, "storageClassName", "<none>")}
}

# Rule ID: storage-class
# All fleet data lives on the encrypted Longhorn storage class
# (policies/policy.yaml). A PVC or volumeClaimTemplate on any other class
# — or with no class, which falls back to the cluster default — is an
# unencrypted-data bug.
deny contains msg if {
	some e in storage_claims
	not e.class == data.storage_class
	msg := sprintf("storage claim '%v' in '%v' requests storage class %v; fleet data must use '%v' (encrypted)", [e.name, input.metadata.namespace, e.class, data.storage_class])
}

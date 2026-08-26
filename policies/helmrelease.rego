package main

import rego.v1

# Rule ID: configmap-hash
# A HelmRelease that takes Helm values from a ConfigMap must carry the
# configMapHash annotation. The hash is stamped from the ConfigMap's data
# (see ComputeConfigMapHash); editing the config changes the annotation,
# which changes the HelmRelease, which is what makes Flux re-reconcile.
# Without it, config edits are silently never applied.
deny contains msg if {
	input.kind == "HelmRelease"
	some vf in input.spec.valuesFrom
	vf.kind == "ConfigMap"
	not input.metadata.annotations.configMapHash
	msg := sprintf("HelmRelease '%v' in '%v' takes values from a ConfigMap but has no configMapHash annotation; config changes will not trigger reconciliation", [input.metadata.name, input.metadata.namespace])
}

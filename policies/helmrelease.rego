package main

import rego.v1

# Rule ID: values-watch
# A ConfigMap holding Helm values (key values.yaml, as CreateHelmValuesConfig
# writes them) must carry the reconcile.fluxcd.io/watch=Enabled label.
# helm-controller only reacts to changes of a HelmRelease's spec, and the
# values ConfigMap keeps its name across edits: without the label, an edit
# waits for the release's next interval (10m) to be applied.
deny contains msg if {
	input.kind == "ConfigMap"
	input.data["values.yaml"]
	not input.metadata.labels["reconcile.fluxcd.io/watch"] == "Enabled"
	msg := sprintf("ConfigMap '%v' in '%v' holds Helm values but has no reconcile.fluxcd.io/watch=Enabled label; helm-controller will only apply edits at its next interval", [input.metadata.name, input.metadata.namespace])
}

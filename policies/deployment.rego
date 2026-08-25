package main

import rego.v1

# Rule ID: run-as-root
# Deployments must declare pod-level runAsNonRoot.
# Exemptions: policies/exemptions.yaml (rule "run-as-root").
deny contains msg if {
	input.kind == "Deployment"
	not input.spec.template.spec.securityContext.runAsNonRoot
	not exempt("run-as-root")
	msg := "Containers must not run as root"
}

# Rule ID: namespace
# Every Deployment must live in a namespace.
deny contains msg if {
	input.kind == "Deployment"
	not input.metadata.namespace
	msg := "Missing namespace"
}

# exempt is true when the resource's namespace has the rule in its
# exemptions.yaml entry. Keyed by namespace because one chart owns exactly
# one namespace (namespace name = app name), while generated resource names
# carry a spec hash that changes when the spec changes.
# Missing entries or missing rules simply do not exempt.
exempt(rule) if {
	entry := object.get(data.exemptions, input.metadata.namespace, {})
	some r in entry.rules
	r == rule
}

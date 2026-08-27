package main

import rego.v1

# Rule ID: no-nodeport
# Services must not open direct node/cluster-level exposure: the types in
# policies/policy.yaml (denied_service_types) are forbidden — exposure goes
# through Traefik (see AGENTS.md). ExternalName services are DNS aliases
# for external hosts and expose no workload, so they are out of scope.
# A legitimate need for a denied type gets a reasoned namespace entry in
# policies/exemptions.yaml.
deny contains msg if {
	input.kind == "Service"
	some t in data.denied_service_types
	input.spec.type == t
	msg := sprintf("Service '%v' in '%v' is type %v; node/cluster-level exposure is forbidden, route through Traefik (see AGENTS.md)", [input.metadata.name, input.metadata.namespace, t])
}

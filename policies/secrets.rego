package main

import rego.v1

# Rule ID: inline-secret
# Secrets belong in Vault + External Secrets Operator (see AGENTS.md).
# A kind: Secret in dist/ is only acceptable where an upstream chart
# requires it (chart-wired config files), and must be listed in
# policies/exemptions.yaml with a reason. This stops plaintext secret
# sprawl from creeping into the repo and the generated branch.
deny contains msg if {
	input.kind == "Secret"
	has_inline_data
	not exempt("inline-secret")
	msg := sprintf("Secret '%v' in '%v' carries inline data; secrets belong in Vault — add a reasoned exemption to policies/exemptions.yaml only if the chart requires it", [input.metadata.name, input.metadata.namespace])
}

# has_inline_data is true when the Secret actually carries data (an empty
# or data-less Secret is not secret sprawl).
has_inline_data if {
	input.data
}

has_inline_data if {
	input.stringData
}

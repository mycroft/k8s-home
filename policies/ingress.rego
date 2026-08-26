package main

import rego.v1

# Rule ID: ingress-tls
# Every Ingress must serve TLS. Plain-HTTP exposure is not a thing in this
# fleet (see AGENTS.md, Ingress).
deny contains msg if {
	input.kind == "Ingress"
	not input.spec.tls
	msg := sprintf("Ingress '%v' in '%v' has no tls block; it would be served over plain HTTP", [input.metadata.name, input.metadata.namespace])
}

# Rule ID: ingress-issuer
# Every Ingress must be issued by the fleet's cert-manager cluster issuer
# (AGENTS.md pins letsencrypt-prod). Without the annotation no certificate
# is ever issued and the TLS block is dead weight.
deny contains msg if {
	input.kind == "Ingress"
	not input.metadata.annotations["cert-manager.io/cluster-issuer"] == "letsencrypt-prod"
	msg := sprintf("Ingress '%v' in '%v' does not use cluster issuer 'letsencrypt-prod'; no certificate will be issued", [input.metadata.name, input.metadata.namespace])
}

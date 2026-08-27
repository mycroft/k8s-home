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

# Rule ID: ingress-domain
# Ingress hosts must belong to a fleet domain (policies/policy.yaml).
# Catches typos — a misspelled host still gets a wildcard cert and DNS —
# and rogue public-facing hostnames.
# The match is on a label boundary: not-iop.cx does not satisfy iop.cx.
deny contains msg if {
	input.kind == "Ingress"
	some r in input.spec.rules
	not allowed_domain(r.host)
	msg := sprintf("Ingress '%v' in '%v' uses host '%v'; allowed suffixes: %v", [input.metadata.name, input.metadata.namespace, r.host, data.allowed_domain_suffixes])
}

allowed_domain(host) if {
	some s in data.allowed_domain_suffixes
	host == s
}

allowed_domain(host) if {
	some s in data.allowed_domain_suffixes
	endswith(host, sprintf(".%s", [s]))
}

package main

import rego.v1

# Rule ID: latest
# Images must not use the mutable :latest tag.
# Exemptions: policies/exemptions.yaml (rule "latest").
deny contains msg if {
	some c in input.spec.template.spec.containers
	endswith(c.image, ":latest")
	not exempt("latest")
	msg := sprintf("image '%v' : latest tag is forbidden", [c.image])
}

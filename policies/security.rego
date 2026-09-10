package main

import rego.v1

# privileged_containers is the set of container names conftest can see in
# dist/ that ask for privileged: true — containers and initContainers of
# Deployments/StatefulSets/DaemonSets (spec.template) and of CronJobs
# (spec.jobTemplate). This mirrors container_images in images.rego, so a
# new workload kind only needs a new line here.
#
# Only container-level securityContext is checked. `privileged` is a field
# of core/v1 SecurityContext, not PodSecurityContext, so there is no
# pod-level equivalent to match.
privileged_containers contains name if {
	c := input.spec.template.spec.containers[_]
	c.securityContext.privileged == true
	name := c.name
}

privileged_containers contains name if {
	c := input.spec.template.spec.initContainers[_]
	c.securityContext.privileged == true
	name := c.name
}

privileged_containers contains name if {
	c := input.spec.jobTemplate.spec.template.spec.containers[_]
	c.securityContext.privileged == true
	name := c.name
}

privileged_containers contains name if {
	c := input.spec.jobTemplate.spec.template.spec.initContainers[_]
	c.securityContext.privileged == true
	name := c.name
}

# Rule ID: privileged
# A privileged container switches off the container boundary: it gets all
# capabilities, host devices, and no seccomp/AppArmor confinement, which
# makes escape to the node trivial. Homelab apps do not need it.
# Exemptions: policies/exemptions.yaml (rule "privileged").
deny contains msg if {
	some name in privileged_containers
	not exempt("privileged")
	msg := sprintf("container '%v' must not run privileged", [name])
}

package main

import rego.v1

# container_images is the set of all image references conftest can see in
# dist/: containers and initContainers of Deployments/StatefulSets
# (spec.template) and of CronJobs (spec.jobTemplate). Both image rules
# below iterate this set, so a new workload kind only needs a new line here.
container_images contains img if {
	img := input.spec.template.spec.containers[_].image
}

container_images contains img if {
	img := input.spec.template.spec.initContainers[_].image
}

container_images contains img if {
	img := input.spec.jobTemplate.spec.template.spec.containers[_].image
}

container_images contains img if {
	img := input.spec.jobTemplate.spec.template.spec.initContainers[_].image
}

# Rule ID: latest
# Images must not use the mutable :latest tag.
# Exemptions: policies/exemptions.yaml (rule "latest").
deny contains msg if {
	some img in container_images
	endswith(img, ":latest")
	not exempt("latest")
	msg := sprintf("image '%v' : latest tag is forbidden", [img])
}

# Rule ID: image-tag
# Every image must carry an explicit tag or digest. A bare
# `image: org/app` resolves to :latest on the cluster — which is exactly
# what a missing versions.yaml entry produces, because
# RegisterContainerImage silently returns the untagged name. Catch it at
# CI time instead of at pod start.
# The tag check looks at the last path segment only, so registry hosts with
# ports (host:5000/app) are not mistaken for tagged images.
deny contains msg if {
	some img in container_images
	not contains(img, "@")
	last_segment := split(img, "/")
	last_part := last_segment[count(last_segment) - 1]
	not contains(last_part, ":")
	msg := sprintf("image '%v' has no explicit tag or digest; Kubernetes resolves it to :latest", [img])
}

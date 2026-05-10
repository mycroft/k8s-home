package main

import rego.v1

warn contains msg if {
    some c in input.spec.template.spec.containers
    endswith(c.image, ":latest")
    msg := sprintf("image '%v' : latest tag is forbidden", [c.image])
}

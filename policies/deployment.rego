package main

warn contains msg if {
  input.kind == "Deployment"
  not input.spec.template.spec.securityContext.runAsNonRoot

  msg := "Containers must not run as root"
}

warn contains msg if {
  input.kind == "Deployment"
  not input.metadata.namespace

  msg := "Missing namespace"
}

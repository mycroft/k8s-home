---
name: versions
description: Version bumps in versions.yaml, which no CI job checks
applies_to: ["versions.yaml"]
---
`versions.yaml` pins every Helm chart and container image deployed to the
cluster. A change to it alone runs no CI and deploys on merge, so this review
is its only check. Most changes are automated bumps opened by
`k8s-home create-prs`; an ordinary patch or minor bump is not a finding.
Report only:

- A release candidate, alpha, beta or other prerelease version: these are
  never deployed here.
- A downgrade.
- A major version change, so that its upgrade notes get read before the
  merge: report it as `minor`, naming the old and new versions.
- A tag that changes scheme: a `v` prefix added or dropped, a variant suffix
  such as `-alpine` added or dropped, a date or hash replacing a semantic
  version.
- An entry under `patterns:` that is not a valid regular expression, or that
  would now match prerelease or variant tags.
- An entry removed while the chart that uses it may still exist: a Helm
  chart missing from this file fails the deploy.

---
name: versions
description: Version bumps in versions.yaml, checked against what changed upstream
applies_to: ["versions.yaml"]
---
`versions.yaml` pins every Helm chart and container image deployed to the
cluster. A change to it alone runs no CI and deploys on merge, and an approved
one is merged automatically, so this review is its only check. Most changes
are automated bumps opened by `k8s-home create-prs`.

For each changed entry, read what changed upstream between the old and the
new version, then check it against how this repository uses it:

- **Helm charts** (`helmcharts:`, written `<repository>/<chart>`): find the
  repository's URL in `charts/` (grep the repository name, e.g.
  `repoName := "hashicorp"` next to `https://helm.releases.hashicorp.com`, or
  `homelab-oci`, declared in `charts/infra/homelab-oci.go`). Call
  `chart_diff` from the old to the new version, then compare what moved in
  `values.yaml` and the templates with the values this repository sets: the
  Go chart that installs it and its `configs/<release>.yaml`.
- **Container images** (`images:`): call `release_notes` for the image
  between the two tags. Look for breaking changes, removed or renamed
  settings and environment variables, and required migrations, and check
  them against the chart in `charts/` that runs the image. Follow a release
  note to its upgrade guide with `fetch_url` when it points to one.
- **Self-built charts and images** (`homelab-oci/*`, `registry.mkz.me/...`):
  call `repo_log` on the source repository on git.mkz.me between the two
  versions; the chart's Go file links to it.

Report:

- `major`: an upstream change that breaks this repository as configured: a
  value or setting it uses that was renamed or removed, a default it relies
  on that changed, a required manual migration (database, storage, CRDs).
  Name the upstream evidence (chart file and line, release, page) and the
  file of this repository it affects.
- `minor`: a major version change, so a person reads the upgrade notes
  before it merges, naming the old and new versions and what you found in
  them; and a bump whose upstream changes you could not read (no chart
  version, no release notes, a tool error), saying what you tried.
- `minor`: a release candidate, alpha, beta or other prerelease version:
  these are never deployed here.
- `minor`: a downgrade; a tag that changes scheme (a `v` prefix added or
  dropped, a variant suffix such as `-alpine` added or dropped, a date or
  hash replacing a semantic version); an entry under `patterns:` that is not
  a valid regular expression or would now match prerelease or variant
  tags; an entry removed while the chart that uses it may still exist.

An ordinary patch or minor bump whose changes do not affect this repository
is not a finding: say what changed in the summary.

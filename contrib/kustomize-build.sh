#!/bin/sh
# Builds dist/ the way Flux builds the generated branch, so a manifest set
# that kustomize-controller would reject fails here instead of after merge
# (e.g. two charts declaring the same HelmRepository, see #1727).
#
# deploy.yaml replaces generated/*.yaml with dist/*, and generated/ has no
# kustomization.yaml, so kustomize-controller generates one: every top-level
# .yaml file plus the flux-system/ bootstrap directory, which has its own.
# This rebuilds that layout from dist/ and the generated branch's
# flux-system/.

set -eu

# CI pins kustomize in lint.yaml to kustomize-controller's kustomize/api. A
# local kustomize, or the one embedded in kubectl, may be another version, so
# the version used is printed with the result.
if command -v kustomize >/dev/null 2>&1; then
	build="kustomize build"
	version=$(kustomize version)
else
	build="kubectl kustomize"
	version=$(kubectl version --client | sed -n 's/^Kustomize Version: //p')
fi

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
mkdir "$work/root"

cp dist/*.yaml "$work/root/"

if ! git rev-parse --verify -q origin/generated >/dev/null; then
	git fetch -q origin generated:refs/remotes/origin/generated
fi
git archive origin/generated generated/flux-system | tar -x -C "$work/root" --strip-components=1

{
	echo "apiVersion: kustomize.config.k8s.io/v1beta1"
	echo "kind: Kustomization"
	echo "resources:"
	echo "  - flux-system"
	for f in dist/*.yaml; do
		echo "  - $(basename "$f")"
	done
} >"$work/root/kustomization.yaml"

$build "$work/root" >"$work/out.yaml"
# A plain assignment, so that an empty build (grep -c finds 0) fails set -e.
count=$(grep -c '^kind:' "$work/out.yaml")
echo "kustomize build OK ($build $version): $count resources"

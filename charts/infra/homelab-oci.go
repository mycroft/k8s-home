package infra

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

// NewHomelabOCIHelmRepositoryChart declares the OCI repository holding the
// Helm charts of self-built apps (fluxcd-ui, agent-nexus). Charts using it
// reference it by name and must not declare it again.
func NewHomelabOCIHelmRepositoryChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	chart := builder.NewChart("homelab-oci-helm-repository")

	chart.CreateHelmRepository(
		"homelab-oci",
		"oci://registry.mkz.me/mycroft/charts",
	)

	return chart
}

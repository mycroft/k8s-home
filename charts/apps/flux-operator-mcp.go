package apps

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

// NewFluxOperatorMCPChart deploys flux-operator-mcp, the Flux MCP server
// (https://github.com/controlplaneio-fluxcd/flux-operator), through which
// agent-nexus reads the state of Flux and of what it deploys. Values are in
// configs/flux-operator-mcp.yaml.
//
// The server runs read-only, and its ServiceAccount, created by the chart,
// reads everything but Secrets instead of the chart's cluster-admin. It has no
// authentication of its own: it has no ingress, and the chart's NetworkPolicy
// only lets agent-nexus in.
func NewFluxOperatorMCPChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "flux-operator-mcp"
	repositoryName := "controlplaneio-fluxcd"
	chartName := "flux-operator-mcp"
	releaseName := "flux-operator-mcp"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repositoryName,
		"oci://ghcr.io/controlplaneio-fluxcd/charts",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	chart.NewReadAllButSecretsClusterRole(releaseName) // the chart names its ServiceAccount after the release

	return chart
}

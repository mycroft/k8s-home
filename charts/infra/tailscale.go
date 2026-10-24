package infra

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

// NewTailscaleChart deploys the Tailscale Kubernetes operator, which serves
// Ingresses of class "tailscale" on the tailnet only (see
// Chart.NewTailscaleIngress and docs/tailscale.md).
func NewTailscaleChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "tailscale"
	repositoryName := "tailscale"
	chartName := "tailscale-operator"
	releaseName := "tailscale-operator"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	// The chart reads its OAuth client from a Secret named operator-oauth
	// (keys client_id and client_secret) when oauth.* is left unset.
	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "operator-oauth")

	chart.CreateHelmRepository(
		repositoryName,
		"https://pkgs.tailscale.com/helmcharts",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	return chart
}

package security

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewVaultChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "vault"
	repoName := "hashicorp"
	releaseName := "vault"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repoName,
		"https://helm.releases.hashicorp.com",
	)

	chart.CreateHelmRelease(
		namespace,
		repoName,
		"vault",     // chart name
		releaseName, // release name
		kubehelpers.WithDefaultConfigFile(),
	)

	// The chart's own ingress (configs/vault.yaml) serves
	// vault.services.mkz.me; this serves https://vault.merino-tiaki.ts.net on
	// the tailnet, in front of the chart's vault Service.
	chart.NewTailscaleIngress("vault", "vault", 8200)

	return chart
}

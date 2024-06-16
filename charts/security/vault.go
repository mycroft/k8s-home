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
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName, // release name to be modified
				"vault.yaml",
			),
		},
		nil,
	)

	return chart
}

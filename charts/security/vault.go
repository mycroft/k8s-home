package security

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewVaultChart(builder *kubehelpers.Builder) cdk8s.Chart {
	namespace := "vault"
	repoName := "hashicorp"
	releaseName := "vault"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateHelmRepository(
		chart.Cdk8sChart,
		repoName,
		"https://helm.releases.hashicorp.com",
	)

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,
		repoName,
		"vault",     // chart name
		releaseName, // release name
		map[string]string{},
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

	return chart.Cdk8sChart
}

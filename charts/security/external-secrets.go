package security

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewExternalSecretsChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "external-secrets"

	repositoryName := "external-secrets"
	chartName := "external-secrets"
	releaseName := "external-secrets"

	chart := builder.NewChart(namespace)

	chart.NewNamespace(namespace)

	kubehelpers.CreateHelmRepository(
		chart.Cdk8sChart,
		repositoryName,
		"https://charts.external-secrets.io",
	)

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,
		repositoryName, // repo name
		chartName,      // chart name
		releaseName,    // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"external-secrets.yaml",
			),
		},
		nil,
	)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "testaroo")

	return chart
}

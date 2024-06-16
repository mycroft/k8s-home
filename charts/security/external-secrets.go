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

	chart.CreateHelmRepository(
		repositoryName,
		"https://charts.external-secrets.io",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName, // repo name
		chartName,      // chart name
		releaseName,    // release name
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

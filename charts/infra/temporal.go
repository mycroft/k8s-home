package infra

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewTemporalChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "temporal"
	repositoryName := "temporal"
	chartName := "temporal"
	releaseName := "temporal"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql-visibility")

	kubehelpers.CreateHelmRepository(
		chart.Cdk8sChart,
		repositoryName,
		"https://go.temporal.io/helm-charts",
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
				"temporal.yaml",
			),
		},
		nil,
	)

	return chart
}

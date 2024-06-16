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

	chart.CreateHelmRepository(
		repositoryName,
		"https://go.temporal.io/helm-charts",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName, // repo name
		chartName,      // chart name
		releaseName,    // release name
		kubehelpers.WithDefaultConfigFile(),
	)

	return chart
}

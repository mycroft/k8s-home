package observability

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewTempoChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "tempo"
	repositoryName := "grafana"
	chartName := "tempo"
	releaseName := "tempo"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"tempo.yaml",
			),
		},
		nil,
	)

	return chart
}

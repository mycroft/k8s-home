package observability

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewPromtailChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "promtail"
	repositoryName := "grafana"
	chartName := "promtail"
	releaseName := "promtail"

	chart := builder.NewChart(namespace)

	chart.NewNamespace(namespace)
	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)

	chart.CreateHelmRelease(
		namespace,
		repositoryName, // repo name; was installed in Loki
		chartName,      // chart name
		chartName,      // release name
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName, // release name to be modified
				"promtail.yaml",
			),
		},
		nil,
	)

	return chart
}

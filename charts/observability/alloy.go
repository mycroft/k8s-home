package observability

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewAlloyChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "alloy"
	repositoryName := "grafana"
	chartName := "alloy"
	releaseName := "alloy"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,
		repositoryName,
		chartName,
		chartName,
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"alloy.yaml",
			),
		},
		nil,
	)

	return chart
}

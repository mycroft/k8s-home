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

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	return chart
}

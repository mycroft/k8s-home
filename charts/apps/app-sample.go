package apps

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewAppSampleChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "app-sample"

	chartName := "app"
	releaseName := namespace
	repositoryName := "mycroft"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repositoryName,
		"https://mycroft.github.io/helm-charts",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	return chart
}

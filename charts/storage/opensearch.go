package storage

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewOpenSearchChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "opensearch"
	namespace := "opensearch"

	chart := builder.NewChart(appName)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		"opensearch",
		"https://opensearch-project.github.io/helm-charts",
	)

	chart.CreateHelmRelease(
		namespace,
		"opensearch", // repoName; must be in flux-system
		"opensearch", // chart name
		"opensearch", // release name
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				"opensearch", // release name
				"opensearch.yaml",
			),
		},
		nil,
	)

	chart.CreateHelmRelease(
		namespace,
		"opensearch",            // repoName; must be in flux-system
		"opensearch-dashboards", // chart name
		"opensearch-dashboards", // release name
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				"opensearch-dashboards", // release name
				"opensearch-dashboards.yaml",
			),
		},
		nil,
	)

	return chart
}

package storage

import "git.mkz.me/mycroft/k8s-home/internal/kubehelpers"

func NewGarage(builder *kubehelpers.Builder) *kubehelpers.Chart {
	chartName := "garage"
	namespace := chartName

	repoName := chartName
	releaseName := chartName

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repoName,
		"https://charts.derwitt.dev/",
	)

	configMaps := []kubehelpers.HelmReleaseConfigMap{
		kubehelpers.CreateHelmValuesConfig(
			chart.Cdk8sChart,
			namespace,
			releaseName,
			"garage.yaml",
		),
	}

	chart.CreateHelmRelease(
		namespace,   // namespace
		repoName,    // repository name, same as above
		chartName,   // the chart name
		releaseName, // the release name
		kubehelpers.WithConfigMaps(configMaps),
	)

	return chart
}

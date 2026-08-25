package apps

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewHomepageChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "homepage"
	namespace := "homepage"
	repositoryName := "jameswynn"
	chartName := "homepage"
	releaseName := "homepage"

	chart := builder.NewChart(appName)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repositoryName,
		"https://jameswynn.github.io/helm-charts",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithConfigMaps(
			[]kubehelpers.HelmReleaseConfigMap{
				chart.CreateHelmValuesConfig(
					namespace,
					releaseName,
					"homepage.yaml",
					nil,
				),
			},
		),
	)

	return chart
}

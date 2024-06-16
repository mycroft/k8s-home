package observability

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewKarmaChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "monitoring"
	appName := "karma"
	repositoryName := "wiremind"
	chartName := "karma"
	releaseName := "karma"

	chart := builder.NewChart(appName)

	chart.CreateHelmRepository(
		repositoryName,
		"https://wiremind.github.io/wiremind-helm-charts",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName, // repoName; must be in flux-system
		chartName,      // chart name
		releaseName,    // release name
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"karma.yaml",
			),
		},
		nil,
	)

	return chart
}

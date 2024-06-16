package apps

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewJitsiChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "jitsi"
	repositoryName := "jitsi"
	chartName := "jitsi-meet"
	releaseName := chartName

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repositoryName,
		"https://jitsi-contrib.github.io/jitsi-helm/",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				repositoryName,
				"jitsi-meet.yaml",
			),
		},
		nil,
	)

	return chart
}

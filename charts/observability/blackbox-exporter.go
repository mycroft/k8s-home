package observability

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewBlackboxExporterChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "blackbox-exporter"
	namespace := "monitoring"
	repositoryName := "prometheus-community"
	chartName := "prometheus-blackbox-exporter"
	releaseName := appName

	chart := builder.NewChart(appName)

	// TODO: Namespace is created in NewKubePrometheusStackChart, as well as HelmRepository
	// This should be refactored

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,
		repositoryName,
		chartName,
		releaseName,
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"blackbox-exporter.yaml",
			),
		},
		nil,
	)

	return chart
}

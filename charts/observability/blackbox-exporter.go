package observability

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewBlackboxExporterChart(scope constructs.Construct) cdk8s.Chart {
	appName := "blackbox-exporter"
	namespace := "monitoring"
	repositoryName := "prometheus-community"
	chartName := "prometheus-blackbox-exporter"
	releaseName := appName

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	// TODO: Namespace is created in NewKubePrometheusStackChart, as well as HelmRepository
	// This should be refactored

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName,
		chartName,
		releaseName,
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"blackbox-exporter.yaml",
			),
		},
		nil,
	)

	return chart
}

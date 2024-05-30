package apps

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewJitsiChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "jitsi"
	repositoryName := "jitsi"
	chartName := "jitsi-meet"
	releaseName := chartName

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://jitsi-contrib.github.io/jitsi-helm/",
	)

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
				repositoryName,
				"jitsi-meet.yaml",
			),
		},
		nil,
	)

	return chart
}

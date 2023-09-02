package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
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

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://jitsi-contrib.github.io/jitsi-helm/",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName,
		chartName,
		releaseName,
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
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

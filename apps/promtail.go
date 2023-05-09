package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewPromtailChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "promtail"
	repositoryName := "grafana"
	chartName := "promtail"
	releaseName := "promtail"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName, // repo name; was installed in Loki
		chartName,      // chart name
		chartName,      // release name
		"6.11.1",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName, // release name to be modified
				"promtail.yaml",
			),
		},
		nil,
	)

	return chart
}

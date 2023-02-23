package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewPromtailChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "promtail"

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
		"grafana",  // repo name; was installed in Loki
		"promtail", // chart name
		"promtail", // release name
		"6.9.0",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"promtail.yaml",
			),
		},
		nil,
	)

	return chart
}

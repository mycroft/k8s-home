package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func CreateNATSChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "nats"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"nats",
		"https://nats-io.github.io/k8s/helm/charts/",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"nats",
		"nats",
		"nats",
		"0.19.9",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"nats.yaml",
			),
		},
		nil,
	)

	return chart
}

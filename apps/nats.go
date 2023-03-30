package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewNATSChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "nats"
	repositoryName := "nats"
	chartName := "nats"
	releaseName := "nats"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://nats-io.github.io/k8s/helm/charts/",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName,
		chartName,
		releaseName,
		"0.19.12",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"nats.yaml",
			),
		},
		nil,
	)

	return chart
}

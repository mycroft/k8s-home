package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewLokiChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "loki"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"grafana",
		"https://grafana.github.io/helm-charts",
	)

	k8s_helpers.CreateExternalSecret(chart, namespace, "minio")

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"grafana", // repo name
		"loki",    // chart name
		"loki",    // release name
		"4.8.0",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"loki.yaml",
			),
		},
		nil,
	)

	return chart
}

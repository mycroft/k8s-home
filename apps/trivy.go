package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewTrivyChart(scope constructs.Construct) cdk8s.Chart {
	appName := "trivy"
	namespace := "trivy-system"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"aqua",
		"https://aquasecurity.github.io/helm-charts/",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"aqua",
		"trivy-operator",
		"trivy-operator",
		"0.12.1",
		map[string]string{
			"trivy.ignoreUnfixed": "true",
		},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"trivy.yaml",
			),
		},
		nil, // annotations
	)

	return chart
}
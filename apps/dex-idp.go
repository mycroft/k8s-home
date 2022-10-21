package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewDexIdpChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "dex-idp"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"dex",
		"https://charts.dexidp.io",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"dex", // repo name
		"dex", // chart name
		"dex", // release name
		"0.12.1",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"dex-idp.yaml",
			),
		},
		nil,
	)

	return chart
}

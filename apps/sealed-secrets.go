package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewSealedSecretsChart(scope constructs.Construct) cdk8s.Chart {
	appName := "sealed-secrets"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, appName)

	k8s_helpers.CreateHelmRepository(
		chart,
		"sealed-secrets",
		"https://bitnami-labs.github.io/sealed-secrets",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		"kube-system",    // namespace
		"sealed-secrets", // repo name
		"sealed-secrets", // chart name
		appName,          // release name
		"2.6.9",
		nil,
		nil,
		nil,
	)

	return chart
}

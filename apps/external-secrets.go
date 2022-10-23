package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewExternalSecretsChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "external-secrets"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"external-secrets",
		"https://charts.external-secrets.io",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"external-secrets", // repo name
		"external-secrets", // chart name
		"external-secrets", // release name
		"0.6.0",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{},
		nil,
	)

	k8s_helpers.CreateSecretStore(chart, namespace)
	k8s_helpers.CreateExternalSecret(chart, namespace, "testaroo")

	return chart
}

package security

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewSealedSecretsChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "sealed-secrets"
	releaseName := namespace

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"sealed-secrets",
		"https://bitnami-labs.github.io/sealed-secrets",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,        // namespace
		"sealed-secrets", // repo name
		"sealed-secrets", // chart name
		releaseName,      // release name
		"2.12.0",
		map[string]string{
			"fullnameOverride": "sealed-secrets-controller",
		},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"sealed-secrets.yaml",
			),
		},
		nil,
	)

	return chart
}

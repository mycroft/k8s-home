package security

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
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

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.CreateHelmRepository(
		chart,
		"sealed-secrets",
		"https://bitnami-labs.github.io/sealed-secrets",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,        // namespace
		"sealed-secrets", // repo name
		"sealed-secrets", // chart name
		releaseName,      // release name
		map[string]string{
			"fullnameOverride": "sealed-secrets-controller",
		},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
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

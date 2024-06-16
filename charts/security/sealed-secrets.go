package security

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewSealedSecretsChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "sealed-secrets"
	releaseName := namespace

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		"sealed-secrets",
		"https://bitnami-labs.github.io/sealed-secrets",
	)

	values := map[string]*string{
		"fullnameOverride": jsii.String("sealed-secrets-controller"),
	}

	chart.CreateHelmRelease(
		namespace,        // namespace
		"sealed-secrets", // repo name
		"sealed-secrets", // chart name
		releaseName,      // release name
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"sealed-secrets.yaml",
			),
		},
		nil,
		kubehelpers.WithValues(values),
	)

	return chart
}

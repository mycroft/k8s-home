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
		kubehelpers.WithDefaultConfigFile(),

		kubehelpers.WithHelmValues(values),
	)

	return chart
}

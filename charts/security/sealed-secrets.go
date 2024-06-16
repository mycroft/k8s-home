package security

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewSealedSecretsChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "sealed-secrets"
	releaseName := namespace

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateHelmRepository(
		chart.Cdk8sChart,
		"sealed-secrets",
		"https://bitnami-labs.github.io/sealed-secrets",
	)

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,        // namespace
		"sealed-secrets", // repo name
		"sealed-secrets", // chart name
		releaseName,      // release name
		map[string]string{
			"fullnameOverride": "sealed-secrets-controller",
		},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"sealed-secrets.yaml",
			),
		},
		nil,
	)

	return chart
}

package security

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewKyvernoChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "kyverno"

	repositoryName := "kyverno"
	chartName := "kyverno"
	releaseName := "kyverno"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repositoryName,
		"https://kyverno.github.io/kyverno/",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName, // repo name
		chartName,      // chart name
		releaseName,    // release name
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"kyverno.yaml",
			),
		},
		nil,
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		"kyverno-policies",
		"kyverno-policies",
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				"kyverno-policies",
				"kyverno-policies.yaml",
			),
		},
		nil,
	)

	return chart
}

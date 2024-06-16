package observability

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewJaegerChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "jaeger"
	repositoryName := "jaegertracing"
	chartName := "jaeger"
	releaseName := chartName

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	chart.CreateHelmRepository(
		repositoryName,
		"https://jaegertracing.github.io/helm-charts",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	return chart
}

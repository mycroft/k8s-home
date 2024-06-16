package observability

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewLokiChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "loki"
	repositoryName := "grafana"
	chartName := "loki"
	releaseName := "loki"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)

	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "minio")

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	return chart
}

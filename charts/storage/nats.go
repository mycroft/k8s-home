package storage

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewNATSChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "nats"
	repositoryName := "nats"
	chartName := "nats"
	releaseName := "nats"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repositoryName,
		"https://nats-io.github.io/k8s/helm/charts/",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"nats.yaml",
			),
		},
		nil,
	)

	return chart
}

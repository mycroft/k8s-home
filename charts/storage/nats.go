package storage

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewNATSChart(builder *kubehelpers.Builder) cdk8s.Chart {
	namespace := "nats"
	repositoryName := "nats"
	chartName := "nats"
	releaseName := "nats"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateHelmRepository(
		chart.Cdk8sChart,
		repositoryName,
		"https://nats-io.github.io/k8s/helm/charts/",
	)

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,
		repositoryName,
		chartName,
		releaseName,
		map[string]string{},
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

	return chart.Cdk8sChart
}

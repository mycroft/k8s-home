package apps

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewN8nChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "n8n"

	chartName := "n8n"
	releaseName := chartName
	repositoryName := "helm-charts"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateSecretStore(namespace)
	chart.CreateExternalSecret(namespace, "postgresql")

	_, redisServiceName := chart.NewRedisStatefulset(namespace)

	chart.CreateHelmRepository(
		repositoryName,
		"https://community-charts.github.io/helm-charts",
	)

	configMaps := []kubehelpers.HelmReleaseConfigMap{
		kubehelpers.CreateHelmValuesTemplatedConfig(
			chart.Cdk8sChart,
			namespace,
			repositoryName,
			"n8n.yaml",
			true,
			map[string]interface{}{
				"Redis": redisServiceName,
			},
		),
	}

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithConfigMaps(configMaps),
	)

	return chart
}

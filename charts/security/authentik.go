package security

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewAuthentikChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "authentik"
	repositoryName := "authentik"
	releaseName := "authentik"
	chartName := "authentik"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "authentik-secret")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "mailrelay")

	_, redisServiceName := kubehelpers.NewRedisStatefulset(chart.Cdk8sChart, namespace)
	_ = fmt.Sprintf("redis://%s:6379", redisServiceName)

	kubehelpers.CreateHelmRepository(
		chart.Cdk8sChart,
		repositoryName,
		"https://charts.goauthentik.io",
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
				"authentik.yaml",
			),
		},
		nil,
	)

	return chart
}

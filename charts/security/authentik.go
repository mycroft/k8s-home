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

	_, redisServiceName := chart.NewRedisStatefulset(namespace)
	_ = fmt.Sprintf("redis://%s:6379", redisServiceName)

	chart.CreateHelmRepository(
		repositoryName,
		"https://charts.goauthentik.io",
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
				"authentik.yaml",
			),
		},
		nil,
	)

	return chart
}

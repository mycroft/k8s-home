package security

import (
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
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql-cnpg")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "authentik-secret")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "mailrelay")

	chart.NewRedisStatefulset(namespace)

	chart.CreateHelmRepository(
		repositoryName,
		"https://charts.goauthentik.io",
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

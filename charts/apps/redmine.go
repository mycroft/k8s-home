package apps

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewRedmineChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "redmine"
	releaseName := namespace

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "mariadb")

	chart.CreateHelmRepository(
		"redmine",
		"oci://registry-1.docker.io/bitnamicharts",
	)

	chart.CreateHelmRelease(
		namespace,   // namespace
		"redmine",   // repo name
		"redmine",   // chart name
		releaseName, // release name
		kubehelpers.WithDefaultConfigFile(),
	)

	return chart
}

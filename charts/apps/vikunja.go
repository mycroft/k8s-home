package apps

import "git.mkz.me/mycroft/k8s-home/internal/kubehelpers"

func NewVikunjaChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "vikunja"

	namespace := appName
	releaseName := namespace
	repoName := namespace
	chartName := namespace

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql")

	chart.CreateHelmRepository(
		"vikunja",
		"oci://kolaente.dev/vikunja",
	)

	chart.CreateHelmRelease(
		namespace,
		repoName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	return chart
}

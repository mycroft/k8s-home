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

	kubehelpers.CreateHelmRepository(
		chart.Cdk8sChart,
		"redmine",
		"oci://registry-1.docker.io/bitnamicharts",
	)

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,   // namespace
		"redmine",   // repo name
		"redmine",   // chart name
		releaseName, // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"redmine.yaml",
			),
		},
		nil,
	)

	return chart
}

package storage

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewMariaDBOperator(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "mariadb-operator"
	releaseName := "mariadb-operator"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateHelmRepository(
		chart.Cdk8sChart,
		"mariadb-operator", // repository name
		"https://mariadb-operator.github.io/mariadb-operator",
	)

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,          // namespace
		"mariadb-operator", // repository name, same as above
		"mariadb-operator", // the chart name
		releaseName,
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"mariadb-operator.yaml",
			),
		},
		nil,
	)

	return chart
}

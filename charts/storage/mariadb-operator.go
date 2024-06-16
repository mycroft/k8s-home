package storage

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewMariaDBOperator(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "mariadb-operator"
	releaseName := "mariadb-operator"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		"mariadb-operator", // repository name
		"https://mariadb-operator.github.io/mariadb-operator",
	)

	chart.CreateHelmRelease(
		namespace,          // namespace
		"mariadb-operator", // repository name, same as above
		"mariadb-operator", // the chart name
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	return chart
}

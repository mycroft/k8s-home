package storage

import "git.mkz.me/mycroft/k8s-home/internal/kubehelpers"

func NewCNPGOperator(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "cnpg-system"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		"cnpg",
		"https://cloudnative-pg.github.io/charts",
	)

	chart.CreateHelmRelease(
		namespace,
		"cnpg",
		"cloudnative-pg",
		"cnpg",
	)

	return chart
}

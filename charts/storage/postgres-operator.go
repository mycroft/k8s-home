package storage

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewPostgresOperator(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "postgres-operator"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		"postgres-operator",
		"https://opensource.zalando.com/postgres-operator/charts/postgres-operator",
	)

	chart.CreateHelmRelease(
		namespace,
		"postgres-operator",
		"postgres-operator",
		"postgres-operator",
		nil,
		nil,
	)

	return chart
}

package storage

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewPostgresOperator(builder *kubehelpers.Builder) cdk8s.Chart {
	namespace := "postgres-operator"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateHelmRepository(
		chart.Cdk8sChart,
		"postgres-operator",
		"https://opensource.zalando.com/postgres-operator/charts/postgres-operator",
	)

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,
		"postgres-operator",
		"postgres-operator",
		"postgres-operator",
		map[string]string{},
		nil,
		nil,
	)

	return chart.Cdk8sChart
}

package storage

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewPostgresOperator(scope constructs.Construct) cdk8s.Chart {
	namespace := "postgres-operator"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.CreateHelmRepository(
		chart,
		"postgres-operator",
		"https://opensource.zalando.com/postgres-operator/charts/postgres-operator",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		"postgres-operator",
		"postgres-operator",
		"postgres-operator",
		map[string]string{},
		nil,
		nil,
	)

	return chart
}

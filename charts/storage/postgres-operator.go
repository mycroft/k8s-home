package storage

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
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

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"postgres-operator",
		"https://opensource.zalando.com/postgres-operator/charts/postgres-operator",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"postgres-operator",
		"postgres-operator",
		"postgres-operator",
		"1.10.0",
		map[string]string{},
		nil,
		nil,
	)

	return chart
}

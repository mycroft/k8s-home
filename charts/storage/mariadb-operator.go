package storage

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewMariaDBOperator(scope constructs.Construct) cdk8s.Chart {
	namespace := "mariadb-operator"
	releaseName := "mariadb-operator"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"mariadb-operator", // repository name
		"https://mariadb-operator.github.io/mariadb-operator",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,          // namespace
		"mariadb-operator", // repository name, same as above
		"mariadb-operator", // the chart name
		releaseName,
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"mariadb-operator.yaml",
			),
		},
		nil,
	)

	return chart
}

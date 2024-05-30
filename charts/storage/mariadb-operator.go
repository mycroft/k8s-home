package storage

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
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

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.CreateHelmRepository(
		chart,
		"mariadb-operator", // repository name
		"https://mariadb-operator.github.io/mariadb-operator",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,          // namespace
		"mariadb-operator", // repository name, same as above
		"mariadb-operator", // the chart name
		releaseName,
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
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

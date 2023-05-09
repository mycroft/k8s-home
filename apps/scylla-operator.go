package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewScyllaOperatorChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "scylla-operator"

	repoName := "scylla"
	releaseName := "scylla-operator"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		repoName,
		"https://scylla-operator-charts.storage.googleapis.com/stable",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		repoName,
		"scylla-operator",
		releaseName,
		"v1.8.1",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"scylla-operator.yaml",
			),
		},
		nil,
	)

	return chart
}

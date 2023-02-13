package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewMinioOperator(scope constructs.Construct) cdk8s.Chart {
	namespace := "minio-operator"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"minio",
		"https://operator.min.io/",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"minio",
		"minio-operator",
		"minio-operator",
		"4.3.7",
		map[string]string{},
		nil,
		nil,
	)

	return chart
}

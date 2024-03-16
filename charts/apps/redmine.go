package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewRedmineChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "redmine"
	releaseName := namespace

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"redmine",
		"oci://registry-1.docker.io/bitnamicharts",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,   // namespace
		"redmine",   // repo name
		"redmine",   // chart name
		releaseName, // release name
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"redmine.yaml",
			),
		},
		nil,
	)

	return chart
}

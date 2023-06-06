package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewKyvernoChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "kyverno"

	repositoryName := "kyverno"
	chartName := "kyverno"
	releaseName := "kyverno"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://kyverno.github.io/kyverno/",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName, // repo name
		chartName,      // chart name
		releaseName,    // release name
		"v3.0.1",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"kyverno.yaml",
			),
		},
		nil,
	)

	return chart
}

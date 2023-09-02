package security

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewTrivyChart(scope constructs.Construct) cdk8s.Chart {
	appName := "trivy"
	namespace := "trivy-system"

	repoName := "aqua"
	releaseName := "trivy-operator"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		repoName,
		"https://aquasecurity.github.io/helm-charts/",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		repoName,
		"trivy-operator",
		releaseName,
		map[string]string{
			"trivy.ignoreUnfixed": "true",
		},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"trivy.yaml",
			),
		},
		nil,
	)

	return chart
}

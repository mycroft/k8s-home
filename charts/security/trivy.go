package security

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
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

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.CreateHelmRepository(
		chart,
		repoName,
		"https://aquasecurity.github.io/helm-charts/",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repoName,
		"trivy-operator",
		releaseName,
		map[string]string{
			"trivy.ignoreUnfixed": "true",
		},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
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

package security

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
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

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://kyverno.github.io/kyverno/",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName, // repo name
		chartName,      // chart name
		releaseName,    // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"kyverno.yaml",
			),
		},
		nil,
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName,
		"kyverno-policies",
		"kyverno-policies",
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"kyverno-policies",
				"kyverno-policies.yaml",
			),
		},
		nil,
	)

	return chart
}

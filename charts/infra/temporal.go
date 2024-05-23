package infra

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewTemporalChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "temporal"
	repositoryName := "temporal"
	chartName := "temporal"
	releaseName := "temporal"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateSecretStore(chart, namespace)
	k8s_helpers.CreateExternalSecret(chart, namespace, "postgresql")
	k8s_helpers.CreateExternalSecret(chart, namespace, "postgresql-visibility")

	k8s_helpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://go.temporal.io/helm-charts",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName, // repo name
		chartName,      // chart name
		releaseName,    // release name
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"temporal.yaml",
			),
		},
		nil,
	)

	return chart
}

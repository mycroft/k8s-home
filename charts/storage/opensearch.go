package storage

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewOpenSearchChart(scope constructs.Construct) cdk8s.Chart {
	appName := "opensearch"
	namespace := "opensearch"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"opensearch",
		"https://opensearch-project.github.io/helm-charts",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"opensearch", // repoName; must be in flux-system
		"opensearch", // chart name
		"opensearch", // release name
		"2.12.0",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"opensearch", // release name
				"opensearch.yaml",
			),
		},
		nil,
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"opensearch",            // repoName; must be in flux-system
		"opensearch-dashboards", // chart name
		"opensearch-dashboards", // release name
		"2.10.0",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"opensearch-dashboards", // release name
				"opensearch-dashboards.yaml",
			),
		},
		nil,
	)

	return chart
}

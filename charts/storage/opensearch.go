package storage

import (
	"context"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewOpenSearchChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	appName := "opensearch"
	namespace := "opensearch"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.CreateHelmRepository(
		chart,
		"opensearch",
		"https://opensearch-project.github.io/helm-charts",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		"opensearch", // repoName; must be in flux-system
		"opensearch", // chart name
		"opensearch", // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"opensearch", // release name
				"opensearch.yaml",
			),
		},
		nil,
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		"opensearch",            // repoName; must be in flux-system
		"opensearch-dashboards", // chart name
		"opensearch-dashboards", // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
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

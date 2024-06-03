package infra

import (
	"context"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewTemporalChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	namespace := "temporal"
	repositoryName := "temporal"
	chartName := "temporal"
	releaseName := "temporal"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.CreateSecretStore(chart, namespace)
	kubehelpers.CreateExternalSecret(chart, namespace, "postgresql")
	kubehelpers.CreateExternalSecret(chart, namespace, "postgresql-visibility")

	kubehelpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://go.temporal.io/helm-charts",
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
				"temporal.yaml",
			),
		},
		nil,
	)

	return chart
}

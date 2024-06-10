package observability

import (
	"context"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewAlloyChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	namespace := "alloy"
	repositoryName := "grafana"
	chartName := "alloy"
	releaseName := "alloy"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName,
		chartName,
		chartName,
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"alloy.yaml",
			),
		},
		nil,
	)

	return chart

}

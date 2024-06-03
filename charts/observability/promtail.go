package observability

import (
	"context"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewPromtailChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	namespace := "promtail"
	repositoryName := "grafana"
	chartName := "promtail"
	releaseName := "promtail"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)
	kubehelpers.CreateSecretStore(chart, namespace)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName, // repo name; was installed in Loki
		chartName,      // chart name
		chartName,      // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName, // release name to be modified
				"promtail.yaml",
			),
		},
		nil,
	)

	return chart
}

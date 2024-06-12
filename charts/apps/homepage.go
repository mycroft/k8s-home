package apps

import (
	"context"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewHomepageChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	appName := "homepage"
	namespace := "homepage"
	repositoryName := "jameswynn"
	chartName := "homepage"
	releaseName := "homepage"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://jameswynn.github.io/helm-charts",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName,
		chartName,
		releaseName,
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesTemplatedConfig(
				chart,
				namespace,
				releaseName,
				"homepage.yaml",
				true,
			),
		},
		nil,
	)

	return chart
}

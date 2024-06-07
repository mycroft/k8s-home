package observability

import (
	"context"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewKarmaChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	namespace := "monitoring"
	appName := "karma"
	repositoryName := "wiremind"
	chartName := "karma"
	releaseName := "karma"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	kubehelpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://wiremind.github.io/wiremind-helm-charts",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName, // repoName; must be in flux-system
		chartName,      // chart name
		releaseName,    // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"karma.yaml",
			),
		},
		nil,
	)

	return chart
}

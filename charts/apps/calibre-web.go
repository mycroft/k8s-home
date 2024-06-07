package apps

// It is required to put an empty database in the container:
// $ cd /config
// $ curl -LO https://github.com/janeczku/calibre-web/raw/master/library/metadata.db

import (
	"context"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewCalibreWebChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	appName := "calibre-web"
	namespace := "calibre-web"
	repositoryName := "calibre-web"
	chartName := "calibre-web"
	releaseName := "calibre-web"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.CreateHelmRepository(
		chart,
		repositoryName,
		"oci://tccr.io/truecharts",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName,
		chartName,
		releaseName,
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"calibre-web.yaml",
			),
		},
		nil,
	)

	return chart
}

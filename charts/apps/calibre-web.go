package apps

// It is required to put an empty database in the container:
// $ cd /config
// $ curl -LO https://github.com/janeczku/calibre-web/raw/master/library/metadata.db

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewCalibreWebChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "calibre-web"
	namespace := "calibre-web"
	repositoryName := "calibre-web"
	chartName := "calibre-web"
	releaseName := "calibre-web"

	chart := builder.NewChart(appName)
	chart.NewNamespace(namespace)

	kubehelpers.CreateHelmRepository(
		chart.Cdk8sChart,
		repositoryName,
		"oci://tccr.io/truecharts",
	)

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,
		repositoryName,
		chartName,
		releaseName,
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"calibre-web.yaml",
			),
		},
		nil,
	)

	return chart
}

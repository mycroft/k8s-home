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

	chart.CreateHelmRepository(
		repositoryName,
		"oci://tccr.io/truecharts",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	return chart
}

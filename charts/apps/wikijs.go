package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/traefikio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewWikiJsChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "wikijs"

	chartName := "wiki"
	releaseName := chartName
	repositoryName := "requarks"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateSecretStore(namespace)
	chart.CreateExternalSecret(namespace, "postgresql")

	traefikio.NewMiddleware(
		chart.Cdk8sChart,
		jsii.String("traefik-vhost-redirect-wikijs"),
		&traefikio.MiddlewareProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("traefik-vhost-redirect-wikijs"),
				Namespace: jsii.String(namespace),
			},
			Spec: &traefikio.MiddlewareSpec{
				RedirectRegex: &traefikio.MiddlewareSpecRedirectRegex{
					Regex:       jsii.String("https://wiki.iop.cx/(.*)"),
					Replacement: jsii.String("https://wiki.services.mkz.me/${1}"),
				},
			},
		},
	)

	chart.CreateHelmRepository(
		repositoryName,
		"https://charts.js.wiki",
	)

	configMaps := []kubehelpers.HelmReleaseConfigMap{
		kubehelpers.CreateHelmValuesConfig(
			chart.Cdk8sChart,
			namespace,
			repositoryName,
			"wikijs.yaml",
		),
	}

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithConfigMaps(configMaps),
	)

	return chart
}

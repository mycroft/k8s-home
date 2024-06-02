package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/traefikio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewWikiJsChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "wikijs"

	chartName := "wiki"
	releaseName := chartName
	repositoryName := "requarks"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)
	kubehelpers.CreateSecretStore(chart, namespace)
	kubehelpers.CreateExternalSecret(chart, namespace, "postgresql")

	traefikio.NewMiddleware(
		chart,
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

	kubehelpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://charts.js.wiki",
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
				repositoryName,
				"wikijs.yaml",
			),
		},
		nil,
	)

	return chart
}

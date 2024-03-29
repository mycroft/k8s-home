package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/traefikcontainous"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
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

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)
	k8s_helpers.CreateExternalSecret(chart, namespace, "postgresql")

	traefikcontainous.NewMiddleware(
		chart,
		jsii.String("traefik-vhost-redirect-wikijs"),
		&traefikcontainous.MiddlewareProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("traefik-vhost-redirect-wikijs"),
				Namespace: jsii.String(namespace),
			},
			Spec: &traefikcontainous.MiddlewareSpec{
				RedirectRegex: &traefikcontainous.MiddlewareSpecRedirectRegex{
					Regex:       jsii.String("https://wiki.iop.cx/(.*)"),
					Replacement: jsii.String("https://wiki.services.mkz.me/${1}"),
				},
			},
		},
	)

	k8s_helpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://charts.js.wiki",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName,
		chartName,
		releaseName,
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
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

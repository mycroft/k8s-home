package apps

import (
	"fmt"
	"strings"

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

	image := chart.Builder.RegisterContainerImage("requarks/wiki")
	if !strings.Contains(image, ":") {
		panic(fmt.Errorf("invalid image repository/tag: Missing ':' in %s", image))
	}

	configMaps := []kubehelpers.HelmReleaseConfigMap{
		kubehelpers.CreateHelmValuesTemplatedConfig(
			chart.Cdk8sChart,
			namespace,
			repositoryName,
			"wikijs.yaml",
			true,
			map[string]interface{}{
				"Image": map[string]string{
					"Repository": strings.Split(image, ":")[0],
					"Tag":        strings.Split(image, ":")[1],
				},
			},
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

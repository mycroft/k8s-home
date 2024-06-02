package security

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/traefikcontainous"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewAuthentikChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "authentik"
	repositoryName := "authentik"
	releaseName := "authentik"
	chartName := "authentik"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)
	kubehelpers.CreateSecretStore(chart, namespace)
	kubehelpers.CreateExternalSecret(chart, namespace, "postgresql")
	kubehelpers.CreateExternalSecret(chart, namespace, "authentik-secret")

	_, redisServiceName := kubehelpers.NewRedisStatefulset(chart, namespace)
	_ = fmt.Sprintf("redis://%s:6379", redisServiceName)

	kubehelpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://charts.goauthentik.io",
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
				"authentik.yaml",
			),
		},
		nil,
	)

	embeddedOutpostUrl := "http://ak-outpost-authentik-embedded-outpost.authentik:9000/outpost.goauthentik.io/auth/traefik"

	traefikcontainous.NewMiddleware(
		chart,
		jsii.String("authentik-forward-auth-middleware"),
		&traefikcontainous.MiddlewareProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("authentik"),
				Namespace: jsii.String(namespace),
			},
			Spec: &traefikcontainous.MiddlewareSpec{
				ForwardAuth: &traefikcontainous.MiddlewareSpecForwardAuth{
					Address:            jsii.String(embeddedOutpostUrl),
					TrustForwardHeader: jsii.Bool(true),
					AuthRequestHeaders: &[]*string{
						jsii.String("X-authentik-username"),
						jsii.String("X-authentik-groups"),
						jsii.String("X-authentik-email"),
						jsii.String("X-authentik-name"),
						jsii.String("X-authentik-uid"),
						jsii.String("X-authentik-jwt"),
						jsii.String("X-authentik-meta-jwks"),
						jsii.String("X-authentik-meta-outpost"),
						jsii.String("X-authentik-meta-provider"),
						jsii.String("X-authentik-meta-app"),
						jsii.String("X-authentik-meta-version"),
					},
				},
			},
		},
	)

	return chart
}

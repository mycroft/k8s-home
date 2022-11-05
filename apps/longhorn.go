package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"

	"git.mkz.me/mycroft/k8s-home/imports/certificates_certmanagerio"
	"git.mkz.me/mycroft/k8s-home/imports/traefikcontainous"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewLonghornChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "longhorn-system"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"longhorn",
		"https://charts.longhorn.io",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"longhorn", // repo name
		"longhorn", // chart name
		"longhorn", // release name
		"1.3.2",
		nil,
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"longhorn.yaml",
			),
		},
		nil,
	)

	k8s_helpers.CreateExternalSecret(chart, namespace, "basic-auth-users")

	certificates_certmanagerio.NewCertificate(
		chart,
		jsii.String("certificate"),
		&certificates_certmanagerio.CertificateProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("secret-tls-www"),
			},
			Spec: &certificates_certmanagerio.CertificateSpec{
				DnsNames: &[]*string{
					jsii.String("longhorn.services.mkz.me"),
				},
				IssuerRef: &certificates_certmanagerio.CertificateSpecIssuerRef{
					Kind: jsii.String("ClusterIssuer"),
					Name: jsii.String("letsencrypt-prod"),
				},
				SecretName: jsii.String("secret-tls-www"),
			},
		},
	)

	traefikcontainous.NewMiddleware(
		chart,
		jsii.String("basic-auth"),
		&traefikcontainous.MiddlewareProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("basic-auth"),
				Namespace: jsii.String(namespace),
			},
			Spec: &traefikcontainous.MiddlewareSpec{
				BasicAuth: &traefikcontainous.MiddlewareSpecBasicAuth{
					Realm:  jsii.String("Longhorn Authentication"),
					Secret: jsii.String("basic-auth-users"),
				},
			},
		},
	)

	traefikcontainous.NewIngressRoute(
		chart,
		jsii.String("ingress-route"),
		&traefikcontainous.IngressRouteProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
			},
			Spec: &traefikcontainous.IngressRouteSpec{
				EntryPoints: &[]*string{
					jsii.String("web"),
					jsii.String("websecure"),
				},
				Routes: &[]*traefikcontainous.IngressRouteSpecRoutes{
					{
						Kind:  traefikcontainous.IngressRouteSpecRoutesKind_RULE,
						Match: jsii.String("Host(`longhorn.services.mkz.me`)"),
						Middlewares: &[]*traefikcontainous.IngressRouteSpecRoutesMiddlewares{
							{
								Name:      jsii.String("basic-auth"),
								Namespace: jsii.String(namespace),
							},
						},
						Services: &[]*traefikcontainous.IngressRouteSpecRoutesServices{
							{
								Kind: traefikcontainous.IngressRouteSpecRoutesServicesKind_SERVICE,
								Name: jsii.String("longhorn-frontend"),
								Port: traefikcontainous.IngressRouteSpecRoutesServicesPort_FromNumber(jsii.Number(80)),
							},
						},
					},
				},
				Tls: &traefikcontainous.IngressRouteSpecTls{
					SecretName: jsii.String("secret-tls-www"),
				},
			},
		},
	)

	return chart
}

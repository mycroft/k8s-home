package infra

import (
	"git.mkz.me/mycroft/k8s-home/imports/certificates_certmanagerio"
	"git.mkz.me/mycroft/k8s-home/imports/traefikio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewTraefikChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "kube-system"
	ingressHost := "traefik.services.mkz.me"
	secretName := "secret-tls-www"

	chart := builder.NewChart("traefik")

	certificates_certmanagerio.NewCertificate(
		chart.Cdk8sChart,
		jsii.String("traefik-dashboard-cert"),
		&certificates_certmanagerio.CertificateProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
			},
			Spec: &certificates_certmanagerio.CertificateSpec{
				SecretName: jsii.String(secretName),
				DnsNames: &[]*string{
					jsii.String(ingressHost),
				},
				IssuerRef: &certificates_certmanagerio.CertificateSpecIssuerRef{
					Name: jsii.String("letsencrypt-prod"),
					Kind: jsii.String("ClusterIssuer"),
				},
			},
		},
	)

	traefikio.NewIngressRoute(
		chart.Cdk8sChart,
		jsii.String("traefik-dashboard"),
		&traefikio.IngressRouteProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("traefik-dashboard"),
				Namespace: jsii.String(namespace),
			},
			Spec: &traefikio.IngressRouteSpec{
				Routes: &[]*traefikio.IngressRouteSpecRoutes{
					{
						Match: jsii.String("Host(`" + ingressHost + "`) && (PathPrefix(`/api`) || PathPrefix(`/dashboard`))"),
						Kind:  traefikio.IngressRouteSpecRoutesKind_RULE,
						Services: &[]*traefikio.IngressRouteSpecRoutesServices{
							{
								Name: jsii.String("api@internal"),
								Kind: traefikio.IngressRouteSpecRoutesServicesKind_TRAEFIK_SERVICE,
							},
						},
						Middlewares: &[]*traefikio.IngressRouteSpecRoutesMiddlewares{
							{
								Name: jsii.String("traefik-forward-auth-traefik-forward-auth@kubernetescrd"),
							},
						},
					},
				},
				Tls: &traefikio.IngressRouteSpecTls{
					SecretName: jsii.String(secretName),
				},
			},
		},
	)

	return chart
}

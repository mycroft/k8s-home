package kubehelpers

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/certificates_certmanagerio"
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/imports/traefikio"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// ExternalApplicationConfig describes an application hosted outside the cluster.
type ExternalApplicationConfig struct {
	Name          string
	Hostname      string
	ExternalName  string
	Port          uint
	CertificateID string
}

// NewExternalApplicationChart exposes an application hosted outside the cluster
// through an ExternalName service and a TLS-enabled Traefik IngressRoute.
func (builder *Builder) NewExternalApplicationChart(config ExternalApplicationConfig) *Chart {
	chart := builder.NewChart(config.Name)
	chart.NewNamespace(config.Name)

	tlsSecretName := fmt.Sprintf("%s-redirect-tls", config.Name)
	certificateID := config.CertificateID
	if certificateID == "" {
		certificateID = fmt.Sprintf("%s-cert", config.Name)
	}

	certificates_certmanagerio.NewCertificate(
		chart.Cdk8sChart,
		jsii.String(certificateID),
		&certificates_certmanagerio.CertificateProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(config.Name),
			},
			Spec: &certificates_certmanagerio.CertificateSpec{
				SecretName: jsii.String(tlsSecretName),
				DnsNames: &[]*string{
					jsii.String(config.Hostname),
				},
				IssuerRef: &certificates_certmanagerio.CertificateSpecIssuerRef{
					Name: jsii.String("letsencrypt-prod"),
					Kind: jsii.String("ClusterIssuer"),
				},
			},
		},
	)

	service := k8s.NewKubeService(
		chart.Cdk8sChart,
		jsii.String("external-api"),
		&k8s.KubeServiceProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(config.Name),
			},
			Spec: &k8s.ServiceSpec{
				Ports: &[]*k8s.ServicePort{
					{
						Name:       jsii.String("http"),
						Port:       jsii.Number(float64(config.Port)),
						TargetPort: k8s.IntOrString_FromNumber(jsii.Number(float64(config.Port))),
					},
				},
				Type:         jsii.String("ExternalName"),
				ExternalName: jsii.String(config.ExternalName),
			},
		},
	)

	traefikio.NewIngressRoute(
		chart.Cdk8sChart,
		jsii.String(fmt.Sprintf("%s-ingress", config.Name)),
		&traefikio.IngressRouteProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(config.Name),
			},
			Spec: &traefikio.IngressRouteSpec{
				EntryPoints: &[]*string{
					jsii.String("web"),
					jsii.String("websecure"),
				},
				Routes: &[]*traefikio.IngressRouteSpecRoutes{
					{
						Kind:  traefikio.IngressRouteSpecRoutesKind_RULE,
						Match: jsii.String(fmt.Sprintf("Host(`%s`)", config.Hostname)),
						Services: &[]*traefikio.IngressRouteSpecRoutesServices{
							{
								Name: service.Name(),
								Port: traefikio.IngressRouteSpecRoutesServicesPort_FromString(
									jsii.String("http"),
								),
							},
						},
					},
				},
				Tls: &traefikio.IngressRouteSpecTls{
					SecretName: jsii.String(tlsSecretName),
				},
			},
		},
	)

	return chart
}

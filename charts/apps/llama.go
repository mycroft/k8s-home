package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/traefikio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"git.mkz.me/mycroft/k8s-home/imports/certificates_certmanagerio"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
)

func NewLlamaChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "llama"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	hostname := "llm-api.iop.cx"

	certificates_certmanagerio.NewCertificate(
	    chart.Cdk8sChart,
	    jsii.String("openjunk-cert"),
	    &certificates_certmanagerio.CertificateProps{
	        Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
			},
	        Spec: &certificates_certmanagerio.CertificateSpec{
	            SecretName: jsii.String("llama-redirect-tls"),
	            DnsNames: &[]*string{
	                jsii.String(hostname),
	            },
	            IssuerRef: &certificates_certmanagerio.CertificateSpecIssuerRef{
	                Name: jsii.String("letsencrypt-prod"),
	                Kind: jsii.String("ClusterIssuer"),
	            },
	        },
	    },
	)

	serviceName := "external-api"

	svc := k8s.NewKubeService(
	    chart.Cdk8sChart,
	    jsii.String(serviceName),
	    &k8s.KubeServiceProps{
	        Metadata: &k8s.ObjectMeta{
	            Namespace: jsii.String(namespace),
	        },
	        Spec: &k8s.ServiceSpec{
	            Ports: &[]*k8s.ServicePort{
	                {
						Name: jsii.String("http"),
	                    Port: jsii.Number(80),
	                    TargetPort: k8s.IntOrString_FromNumber(jsii.Number(80)),
	                },
	            },
	            Type: jsii.String("ExternalName"),
	            ExternalName: jsii.String("10.0.0.7"),
	        },
	    },
	)

	traefikio.NewIngressRoute(
		chart.Cdk8sChart,
		jsii.String("llama-ingress"),
		&traefikio.IngressRouteProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
			},
			Spec: &traefikio.IngressRouteSpec{
				EntryPoints: &[]*string{
					jsii.String("web"),
					jsii.String("websecure"),
				},
				Routes: &[]*traefikio.IngressRouteSpecRoutes{
					{
						Kind:  traefikio.IngressRouteSpecRoutesKind_RULE,
						Match: jsii.String("Host(`llm-api.iop.cx`)"),
						Services: &[]*traefikio.IngressRouteSpecRoutesServices{
							{
								Name: svc.Name(),
								Port: traefikio.IngressRouteSpecRoutesServicesPort_FromString(jsii.String("http")),
							},
						},
					},
				},
				Tls: &traefikio.IngressRouteSpecTls{
					SecretName: jsii.String("llama-redirect-tls"),
				},
			},
		},
	)

	return chart
}

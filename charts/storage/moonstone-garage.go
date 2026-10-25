package storage

import (
	"git.mkz.me/mycroft/k8s-home/imports/certificates_certmanagerio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// NewMoonstoneGarageChart holds the TLS certificate of the Garage S3 API running on moonstone,
// outside the cluster. Its lan.mkz.me name is validated over DNS-01 (see the cert-manager chart).
func NewMoonstoneGarageChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "moonstone-garage"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	certificates_certmanagerio.NewCertificate(
		chart.Cdk8sChart,
		jsii.String("certificate"),
		&certificates_certmanagerio.CertificateProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("s3-tls"),
			},
			Spec: &certificates_certmanagerio.CertificateSpec{
				DnsNames: &[]*string{
					jsii.String("s3.lan.mkz.me"),
				},
				// Staging until the DNS-01 path is proven end to end, then letsencrypt-prod.
				IssuerRef: &certificates_certmanagerio.CertificateSpecIssuerRef{
					Kind: jsii.String("ClusterIssuer"),
					Name: jsii.String("letsencrypt-staging"),
				},
				SecretName: jsii.String("s3-tls"),
			},
		},
	)

	return chart
}

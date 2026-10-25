package storage

import (
	"git.mkz.me/mycroft/k8s-home/imports/certificates_certmanagerio"
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// NewMoonstoneGarageChart holds the TLS certificate of the Garage S3 API running on moonstone,
// outside the cluster. Its lan.mkz.me name is validated over DNS-01 (see the cert-manager chart).
// moonstone pulls the Secret with contrib/moonstone-garage/garage-tls-sync (docs/garage.md).
func NewMoonstoneGarageChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "moonstone-garage"
	secretName := "s3-tls"
	reader := "moonstone"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	certificates_certmanagerio.NewCertificate(
		chart.Cdk8sChart,
		jsii.String("certificate"),
		&certificates_certmanagerio.CertificateProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String(secretName),
			},
			Spec: &certificates_certmanagerio.CertificateSpec{
				DnsNames: &[]*string{
					jsii.String("s3.lan.mkz.me"),
				},
				IssuerRef: &certificates_certmanagerio.CertificateSpecIssuerRef{
					Kind: jsii.String("ClusterIssuer"),
					Name: jsii.String("letsencrypt-prod"),
				},
				SecretName: jsii.String(secretName),
			},
		},
	)

	// moonstone's identity: it may read the certificate Secret and nothing else.
	k8s.NewKubeServiceAccount(
		chart.Cdk8sChart,
		jsii.String("sa"),
		&k8s.KubeServiceAccountProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
				Name:      jsii.String(reader),
			},
		},
	)

	k8s.NewKubeRole(
		chart.Cdk8sChart,
		jsii.String("role"),
		&k8s.KubeRoleProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
				Name:      jsii.String(secretName + "-reader"),
			},
			Rules: &[]*k8s.PolicyRule{
				{
					ApiGroups:     jsii.Strings(""),
					Resources:     jsii.Strings("secrets"),
					ResourceNames: jsii.Strings(secretName),
					Verbs:         jsii.Strings("get"),
				},
			},
		},
	)

	k8s.NewKubeRoleBinding(
		chart.Cdk8sChart,
		jsii.String("rolebinding"),
		&k8s.KubeRoleBindingProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
				Name:      jsii.String(secretName + "-reader"),
			},
			Subjects: &[]*k8s.Subject{
				{
					Kind:      jsii.String("ServiceAccount"),
					Name:      jsii.String(reader),
					Namespace: jsii.String(namespace),
				},
			},
			RoleRef: &k8s.RoleRef{
				Kind:     jsii.String("Role"),
				Name:     jsii.String(secretName + "-reader"),
				ApiGroup: jsii.String("rbac.authorization.k8s.io"),
			},
		},
	)

	// A long-lived token for moonstone, populated by the control plane. Projected tokens
	// expire, and moonstone has no way to renew one outside the cluster.
	k8s.NewKubeSecret(
		chart.Cdk8sChart,
		jsii.String("token"),
		&k8s.KubeSecretProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
				Name:      jsii.String(reader + "-token"),
				Annotations: &map[string]*string{
					"kubernetes.io/service-account.name": jsii.String(reader),
				},
			},
			Type: jsii.String("kubernetes.io/service-account-token"),
		},
	)

	return chart
}

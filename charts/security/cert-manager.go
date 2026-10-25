package security

import (
	"git.mkz.me/mycroft/k8s-home/imports/certmanagerio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func createClusterIssueur(chart constructs.Construct, name, server string) certmanagerio.ClusterIssuer {
	return certmanagerio.NewClusterIssuer(
		chart,
		jsii.String("cluster-issueur-"+name),
		&certmanagerio.ClusterIssuerProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name: jsii.String(name),
			},
			Spec: &certmanagerio.ClusterIssuerSpec{
				Acme: &certmanagerio.ClusterIssuerSpecAcme{
					Email:  jsii.String("pm+letsencrypt@mkz.me"),
					Server: jsii.String(server),
					PrivateKeySecretRef: &certmanagerio.ClusterIssuerSpecAcmePrivateKeySecretRef{
						Name: jsii.String(name),
					},
					Solvers: &[]*certmanagerio.ClusterIssuerSpecAcmeSolvers{
						// lan.mkz.me names resolve to private addresses Let's Encrypt cannot reach,
						// so prove ownership with a TXT record pushed to maki's BIND instead.
						// The TSIG key may only update _acme-challenge names (update-policy).
						{
							Selector: &certmanagerio.ClusterIssuerSpecAcmeSolversSelector{
								DnsZones: &[]*string{jsii.String("lan.mkz.me")},
							},
							Dns01: &certmanagerio.ClusterIssuerSpecAcmeSolversDns01{
								Rfc2136: &certmanagerio.ClusterIssuerSpecAcmeSolversDns01Rfc2136{
									Nameserver:    jsii.String("maki.mkz.me:53"),
									TsigKeyName:   jsii.String("cert-manager"),
									TsigAlgorithm: jsii.String("HMACSHA256"),
									TsigSecretSecretRef: &certmanagerio.ClusterIssuerSpecAcmeSolversDns01Rfc2136TsigSecretSecretRef{
										Name: jsii.String("rfc2136-tsig"),
										Key:  jsii.String("secret"),
									},
								},
							},
						},
						{
							Http01: &certmanagerio.ClusterIssuerSpecAcmeSolversHttp01{
								Ingress: &certmanagerio.ClusterIssuerSpecAcmeSolversHttp01Ingress{
									Class: jsii.String("traefik"),
								},
							},
						},
					},
				},
			},
		},
	)
}

func NewCertManagerChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "cert-manager"
	appName := "cert-manager"

	chart := builder.NewChart(namespace)

	// create a namespace for cert-manager
	// reason to create the namespace is that flux will append the release name using the targetNamespace used.
	// therefore, the HelmRepository will lie in fluxcd, while HelmRelease will live in cert-manager.
	chart.NewNamespace(namespace)
	chart.CreateHelmRepository(
		"jetstack",
		"https://charts.jetstack.io",
	)

	chart.CreateHelmRelease(
		namespace,
		"jetstack", // repository name
		appName,    // chart name
		appName,    // release name
		kubehelpers.WithDefaultConfigFile(),
	)

	// TSIG secret for the DNS-01 solver; ClusterIssuers read secrets from this namespace.
	chart.CreateSecretStore(namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "rfc2136-tsig")

	// flux does not like having those here with the helm creation just before
	// This should be moved in their own chart
	createClusterIssueur(chart.Cdk8sChart, "letsencrypt-staging", "https://acme-staging-v02.api.letsencrypt.org/directory")
	createClusterIssueur(chart.Cdk8sChart, "letsencrypt-prod", "https://acme-v02.api.letsencrypt.org/directory")

	return chart
}

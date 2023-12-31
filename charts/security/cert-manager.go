package security

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/certmanagerio"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func createClusterIssueur(chart constructs.Construct, name, server string) certmanagerio.ClusterIssuer {
	return certmanagerio.NewClusterIssuer(
		chart,
		jsii.String(fmt.Sprintf("cluster-issueur-%s", name)),
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

func NewCertManagerChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "cert-manager"
	appName := "cert-manager"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	// create a namespace for cert-manager
	// reason to create the namespace is that flux will append the release name using the targetNamespace used.
	// therefore, the HelmRepository will lie in fluxcd, while HelmRelease will live in cert-manager.
	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"jetstack",
		"https://charts.jetstack.io",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"jetstack", // repository name
		appName,    // chart name
		appName,    // release name
		map[string]string{
			"installCRDs": "true",
		},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				appName, // release name
				"cert-manager.yaml",
			),
		},
		nil,
	)

	// flux does not like having those here with the helm creation just before
	// This should be moved in their own chart
	createClusterIssueur(chart, "letsencrypt-staging", "https://acme-staging-v02.api.letsencrypt.org/directory")
	createClusterIssueur(chart, "letsencrypt-prod", "https://acme-v02.api.letsencrypt.org/directory")

	return chart
}

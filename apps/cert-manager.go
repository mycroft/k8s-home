package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/certmanagerio"
	"git.mkz.me/mycroft/k8s-home/imports/helmtoolkitfluxcdio"
	"git.mkz.me/mycroft/k8s-home/imports/sourcetoolkitfluxcdio"
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
	appName := "cert-manager"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	// helm repo add jetstack https://charts.jetstack.io

	sourcetoolkitfluxcdio.NewHelmRepository(
		chart,
		jsii.String(fmt.Sprintf("helm-repo-%s", appName)),
		&sourcetoolkitfluxcdio.HelmRepositoryProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("jetstack"),
				Namespace: jsii.String("default"),
			},
			Spec: &sourcetoolkitfluxcdio.HelmRepositorySpec{
				Url:      jsii.String("https://charts.jetstack.io"),
				Interval: jsii.String("1m0s"),
			},
		},
	)

	// helm install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --version v1.9.1 --set installCRDs=true

	// install:
	//   createNamespace: true
	//   installCRDs: true

	helmtoolkitfluxcdio.NewHelmRelease(
		chart,
		jsii.String(fmt.Sprintf("helm-rel-%s", appName)),
		&helmtoolkitfluxcdio.HelmReleaseProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String(appName),
				Namespace: jsii.String("cert-manager"),
			},
			Spec: &helmtoolkitfluxcdio.HelmReleaseSpec{
				Install: &helmtoolkitfluxcdio.HelmReleaseSpecInstall{
					CreateNamespace: jsii.Bool(true),
					SkipCrDs:        jsii.Bool(false),
				},
				Chart: &helmtoolkitfluxcdio.HelmReleaseSpecChart{
					Spec: &helmtoolkitfluxcdio.HelmReleaseSpecChartSpec{
						Chart: jsii.String("cert-manager"),
						SourceRef: &helmtoolkitfluxcdio.HelmReleaseSpecChartSpecSourceRef{
							Kind:      helmtoolkitfluxcdio.HelmReleaseSpecChartSpecSourceRefKind_HELM_REPOSITORY,
							Name:      jsii.String("jetstack"),
							Namespace: jsii.String("default"),
						},
						Version: jsii.String("v1.9.1"),
					},
				},
				Interval: jsii.String("1m0s"),
				Values: map[string]*string{
					"installCRDs": jsii.String("true"),
				},
			},
		},
	)

	createClusterIssueur(chart, "letsencrypt-staging", "https://acme-staging-v02.api.letsencrypt.org/directory")
	createClusterIssueur(chart, "letsencrypt-prod", "https://acme-v02.api.letsencrypt.org/directory")

	return chart
}

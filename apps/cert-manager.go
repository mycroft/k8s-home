package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/certmanagerio"
	"git.mkz.me/mycroft/k8s-home/imports/helmtoolkitfluxcdio"
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
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
	appName := "cert-manager"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	// create a namespace for cert-manager
	// reason to create the namespace is that flux will append the release name using the targetNamespace used.
	// therefore, the HelmRepository will lie in fluxcd, while HelmRelease will live in cert-manager.
	k8s.NewKubeNamespace(
		chart,
		jsii.String(fmt.Sprintf("ns-%s", appName)),
		&k8s.KubeNamespaceProps{
			Metadata: &k8s.ObjectMeta{
				Name: &appName,
			},
		},
	)

	// helm repo add jetstack https://charts.jetstack.io

	k8s_helpers.CreateHelmRepository(chart, "jetstack", "https://charts.jetstack.io")

	// helm install cert-manager jetstack/cert-manager --namespace cert-manager --version v1.9.1 --set installCRDs=true

	helmtoolkitfluxcdio.NewHelmRelease(
		chart,
		jsii.String(fmt.Sprintf("helm-rel-%s", appName)),
		&helmtoolkitfluxcdio.HelmReleaseProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String(appName),
				Namespace: jsii.String(appName),
			},
			Spec: &helmtoolkitfluxcdio.HelmReleaseSpec{
				Install: &helmtoolkitfluxcdio.HelmReleaseSpecInstall{
					CreateNamespace: jsii.Bool(false),
					SkipCrDs:        jsii.Bool(false),
				},
				Chart: &helmtoolkitfluxcdio.HelmReleaseSpecChart{
					Spec: &helmtoolkitfluxcdio.HelmReleaseSpecChartSpec{
						Chart: jsii.String("cert-manager"),
						SourceRef: &helmtoolkitfluxcdio.HelmReleaseSpecChartSpecSourceRef{
							Kind:      helmtoolkitfluxcdio.HelmReleaseSpecChartSpecSourceRefKind_HELM_REPOSITORY,
							Name:      jsii.String("jetstack"),
							Namespace: jsii.String("flux-system"),
						},
						Version: jsii.String("v1.9.1"),
					},
				},
				Interval: jsii.String("1m0s"),
			},
		},
	)

	createClusterIssueur(chart, "letsencrypt-staging", "https://acme-staging-v02.api.letsencrypt.org/directory")
	createClusterIssueur(chart, "letsencrypt-prod", "https://acme-v02.api.letsencrypt.org/directory")

	return chart
}

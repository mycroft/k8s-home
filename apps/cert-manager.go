package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/helmtoolkitfluxcdio"
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/imports/sourcetoolkitfluxcdio"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

/*
apiVersion: source.toolkit.fluxcd.io/v1beta2
kind: HelmRepository
metadata:

	name: traefik
	namespace: traefik

spec:

	interval: 1m0s
	url: https://helm.traefik.io/traefik

---
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:

	name: podinfo
	namespace: default

spec:

	interval: 5m
	chart:
	  spec:
	    chart: <name|path>
	    version: '4.0.x'
	    sourceRef:
	      kind: <HelmRepository|GitRepository|Bucket>
	      name: podinfo
	      namespace: flux-system
	    interval: 1m
	values:
	  replicaCount: 2
*/
func NewCertManagerChart(scope constructs.Construct) cdk8s.Chart {
	appName := "cert-manager"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	k8s.NewKubeNamespace(
		chart,
		jsii.String(appName),
		&k8s.KubeNamespaceProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String(appName),
			},
		},
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

	helmtoolkitfluxcdio.NewHelmRelease(
		chart,
		jsii.String(fmt.Sprintf("helm-rel-%s", appName)),
		&helmtoolkitfluxcdio.HelmReleaseProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String(appName),
				Namespace: jsii.String("cert-manager"),
			},
			Spec: &helmtoolkitfluxcdio.HelmReleaseSpec{
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

	return chart
}

/*
  values:
    installCRDs: true

# /flux/boot/traefik/helmrelease.yaml
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
spec:
  chart:
    spec:
      chart: traefik
      sourceRef:
        kind: HelmRepository
        name: traefik
      version: 9.18.2
  interval: 1m0s
*/

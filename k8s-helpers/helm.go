package k8s_helpers

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/helmtoolkitfluxcdio"
	"git.mkz.me/mycroft/k8s-home/imports/sourcetoolkitfluxcdio"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// CreateHelmRepository creates a HelmRepository into the flux-system namespace
// similar to:
// - helm repo add jetstack https://charts.jetstack.io
func CreateHelmRepository(chart constructs.Construct, name, url string) sourcetoolkitfluxcdio.HelmRepository {
	return sourcetoolkitfluxcdio.NewHelmRepository(
		chart,
		jsii.String(fmt.Sprintf("helm-repo-%s", name)),
		&sourcetoolkitfluxcdio.HelmRepositoryProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String(name),
				Namespace: jsii.String("flux-system"),
			},
			Spec: &sourcetoolkitfluxcdio.HelmRepositorySpec{
				Url:      jsii.String(url),
				Interval: jsii.String("1m0s"),
			},
		},
	)
}

type HelmReleaseConfigMap struct {
	Name    string // ConfigMap name
	KeyName string // key name
}

// CreateHelmRelease creates a helm release in the given namespace for the given repo/name and version
// It installs CRDs by default.
// ex:
// - helm install cert-manager jetstack/cert-manager --namespace cert-manager --version v1.9.1 --set installCRDs=true
func CreateHelmRelease(
	chart constructs.Construct,
	namespace, repoName, chartName, releaseName, version string,
	values map[string]string,
	configMaps []HelmReleaseConfigMap,
	annotations map[string]*string,
) helmtoolkitfluxcdio.HelmRelease {
	// Prepare configMaps.
	valuesFrom := []*helmtoolkitfluxcdio.HelmReleaseSpecValuesFrom{}
	for _, configMap := range configMaps {
		valuesFrom = append(valuesFrom, &helmtoolkitfluxcdio.HelmReleaseSpecValuesFrom{
			Kind:      helmtoolkitfluxcdio.HelmReleaseSpecValuesFromKind_CONFIG_MAP,
			Name:      jsii.String(configMap.Name),
			ValuesKey: jsii.String(configMap.KeyName),
		})
	}

	return helmtoolkitfluxcdio.NewHelmRelease(
		chart,
		jsii.String(fmt.Sprintf("helm-rel-%s", releaseName)),
		&helmtoolkitfluxcdio.HelmReleaseProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:        jsii.String(releaseName),
				Namespace:   jsii.String(namespace),
				Annotations: &annotations,
			},
			Spec: &helmtoolkitfluxcdio.HelmReleaseSpec{
				Install: &helmtoolkitfluxcdio.HelmReleaseSpecInstall{
					CreateNamespace: jsii.Bool(false),
					SkipCrDs:        jsii.Bool(false),
				},
				Chart: &helmtoolkitfluxcdio.HelmReleaseSpecChart{
					Spec: &helmtoolkitfluxcdio.HelmReleaseSpecChartSpec{
						Chart: jsii.String(chartName),
						SourceRef: &helmtoolkitfluxcdio.HelmReleaseSpecChartSpecSourceRef{
							Kind:      helmtoolkitfluxcdio.HelmReleaseSpecChartSpecSourceRefKind_HELM_REPOSITORY,
							Name:      jsii.String(repoName),
							Namespace: jsii.String("flux-system"),
						},
						Version: jsii.String(version),
					},
				},
				Interval:   jsii.String("1m0s"),
				Values:     values,
				ValuesFrom: &valuesFrom,
			},
		},
	)
}

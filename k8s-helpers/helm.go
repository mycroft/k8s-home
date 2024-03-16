package k8s_helpers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.mkz.me/mycroft/k8s-home/imports/helmtoolkitfluxcdio"
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/imports/sourcetoolkitfluxcdio"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type HelmChartVersion struct {
	RepositoryName string
	ChartName      string
	Version        string
}

var helmRepositories = map[string]string{}
var helmChartVersions = []HelmChartVersion{}

// CreateHelmRepository creates a HelmRepository into the flux-system namespace
// similar to:
// - helm repo add jetstack https://charts.jetstack.io
func CreateHelmRepository(chart constructs.Construct, name, url string) sourcetoolkitfluxcdio.HelmRepositoryV1Beta2 {
	helmRepositories[name] = url

	spec := sourcetoolkitfluxcdio.HelmRepositoryV1Beta2Spec{
		Url:      jsii.String(url),
		Interval: jsii.String("1m0s"),
	}

	if strings.HasPrefix(url, "oci://") {
		spec.Type = sourcetoolkitfluxcdio.HelmRepositoryV1Beta2SpecType_OCI
	}

	return sourcetoolkitfluxcdio.NewHelmRepositoryV1Beta2(
		chart,
		jsii.String(fmt.Sprintf("helm-repo-%s", name)),
		&sourcetoolkitfluxcdio.HelmRepositoryV1Beta2Props{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String(name),
				Namespace: jsii.String("flux-system"),
			},
			Spec: &spec,
		},
	)
}

type HelmReleaseConfigMap struct {
	Name          string // ConfigMap name
	KeyName       string // key name
	ConfigMapHash string // The Hash of the configmap content
}

// CreateHelmRelease creates a helm release in the given namespace for the given repo/name and version
// It installs CRDs by default.
// ex:
// - helm install cert-manager jetstack/cert-manager --namespace cert-manager --version v1.9.1 --set installCRDs=true
func CreateHelmRelease(
	chart constructs.Construct,
	namespace, repoName, chartName, releaseName string,
	values map[string]string,
	configMaps []HelmReleaseConfigMap,
	annotations map[string]*string,
) helmtoolkitfluxcdio.HelmRelease {
	cachedVersion := ""

	ReadVersions()

	if version, exists := versions.HelmCharts[fmt.Sprintf("%s/%s", repoName, chartName)]; exists {
		cachedVersion = version
	} else {
		panic(fmt.Sprintf("Could not find version for HelmRelease %s/%s", repoName, chartName))
	}

	helmChartVersions = append(helmChartVersions, HelmChartVersion{
		RepositoryName: repoName,
		ChartName:      chartName,
		Version:        cachedVersion,
	})

	// Prepare configMaps.
	valuesFrom := []*helmtoolkitfluxcdio.HelmReleaseSpecValuesFrom{}
	for _, configMap := range configMaps {
		valuesFrom = append(valuesFrom, &helmtoolkitfluxcdio.HelmReleaseSpecValuesFrom{
			Kind:      helmtoolkitfluxcdio.HelmReleaseSpecValuesFromKind_CONFIG_MAP,
			Name:      jsii.String(configMap.Name),
			ValuesKey: jsii.String(configMap.KeyName),
		})
		if annotations == nil {
			annotations = map[string]*string{}
		}
		annotations["configMapHash"] = jsii.String(configMap.ConfigMapHash)
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
						Version: jsii.String(cachedVersion),
					},
				},
				Interval:   jsii.String("1m0s"),
				Values:     values,
				ValuesFrom: &valuesFrom,
			},
		},
	)
}

func CreateHelmValuesConfig(
	chart constructs.Construct,
	namespace, releaseName, filename string,
) HelmReleaseConfigMap {
	filepath := filepath.Join("configs", filename)
	contents, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	constructName := "helm-values"
	if releaseName != "" {
		constructName = fmt.Sprintf("helm-val-%s", releaseName)
	} else {
		log.Printf("WARNING: HelmValues in ns:%s is still using legacy name", namespace)
	}

	cm := k8s.NewKubeConfigMap(
		chart,
		jsii.String(constructName),
		&k8s.KubeConfigMapProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
			},
			Data: &map[string]*string{
				"values.yaml": jsii.String(string(contents)),
			},
		},
	)

	return HelmReleaseConfigMap{
		Name:          *cm.Name(),
		KeyName:       "values.yaml",
		ConfigMapHash: ComputeConfigMapHash(cm),
	}
}

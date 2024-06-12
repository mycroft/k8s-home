package kubehelpers

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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

type TemplateValues struct {
	Hash string
}

var helmRepositories = map[string]string{}
var helmChartVersions = []HelmChartVersion{}

// CreateHelmRepository creates a HelmRepository into the flux-system namespace
// similar to:
// - helm repo add jetstack https://charts.jetstack.io
func CreateHelmRepository(chart constructs.Construct, name, url string) sourcetoolkitfluxcdio.HelmRepository {
	helmRepositories[name] = url

	spec := sourcetoolkitfluxcdio.HelmRepositorySpec{
		Url:      jsii.String(url),
		Interval: jsii.String("1m0s"),
	}

	if strings.HasPrefix(url, "oci://") {
		spec.Type = sourcetoolkitfluxcdio.HelmRepositorySpecType_OCI
	}

	return sourcetoolkitfluxcdio.NewHelmRepository(
		chart,
		jsii.String(fmt.Sprintf("helm-repo-%s", name)),
		&sourcetoolkitfluxcdio.HelmRepositoryProps{
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
				Interval:   jsii.String("10m0s"),
				Timeout:    jsii.String("5m0s"),
				Values:     values,
				ValuesFrom: &valuesFrom,
			},
		},
	)
}

func CreateHelmValuesTemplatedConfig(
	chart constructs.Construct,
	namespace, releaseName, filename string,
	useCustomTemplate bool,
) HelmReleaseConfigMap {
	var doc bytes.Buffer

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

	renderedContents := string(contents)

	if useCustomTemplate {
		// Apply custom templating, starting with sha256 of config file
		h := sha256.New()
		h.Write([]byte(renderedContents))

		values := TemplateValues{
			Hash: fmt.Sprintf("%x", h.Sum(nil)),
		}

		tmpl, err := template.New("config").Parse(renderedContents)
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(&doc, values)
		if err != nil {
			panic(err)
		}

		renderedContents = doc.String()
	}

	cm := k8s.NewKubeConfigMap(
		chart,
		jsii.String(constructName),
		&k8s.KubeConfigMapProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
			},
			Data: &map[string]*string{
				"values.yaml": jsii.String(renderedContents),
			},
		},
	)

	return HelmReleaseConfigMap{
		Name:          *cm.Name(),
		KeyName:       "values.yaml",
		ConfigMapHash: ComputeConfigMapHash(cm),
	}
}

func CreateHelmValuesConfig(
	chart constructs.Construct,
	namespace, releaseName, filename string,
) HelmReleaseConfigMap {
	return CreateHelmValuesTemplatedConfig(
		chart,
		namespace,
		releaseName,
		filename,
		false,
	)
}

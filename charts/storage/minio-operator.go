package storage

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewMinioOperator(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "minio-operator"

	repoName := "minio"
	releaseName := "operator"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)

	chart.CreateHelmRepository(
		repoName,
		"https://operator.min.io/",
	)

	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "root")

	configMaps := []kubehelpers.HelmReleaseConfigMap{
		kubehelpers.CreateHelmValuesConfig(
			chart.Cdk8sChart,
			namespace,
			releaseName,
			"minio-operator.yaml",
		),
	}

	chart.CreateHelmRelease(
		namespace,   // namespace
		repoName,    // repository name, same as above
		"operator",  // the chart name
		releaseName, // the release name
		kubehelpers.WithConfigMaps(configMaps),
	)

	k8s.NewKubeSecret(
		chart.Cdk8sChart,
		jsii.String("console-sa-secret"),
		&k8s.KubeSecretProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("console-sa-secret"),
				Namespace: jsii.String(namespace),
				Annotations: &map[string]*string{
					"kubernetes.io/service-account.name": jsii.String("console-sa"),
				},
			},
			Type: jsii.String("kubernetes.io/service-account-token"),
		},
	)

	return chart
}

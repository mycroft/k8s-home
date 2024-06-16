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

	kubehelpers.CreateHelmRepository(
		chart.Cdk8sChart,
		repoName,
		"https://operator.min.io/",
	)

	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "root")

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,   // namespace
		repoName,    // repository name, same as above
		"operator",  // the chart name
		releaseName, // the release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"minio-operator.yaml",
			),
		},
		nil,
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

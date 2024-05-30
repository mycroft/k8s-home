package storage

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewMinioOperator(scope constructs.Construct) cdk8s.Chart {
	namespace := "minio-operator"

	repoName := "minio"
	releaseName := "operator"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)
	kubehelpers.CreateSecretStore(chart, namespace)

	kubehelpers.CreateHelmRepository(
		chart,
		repoName,
		"https://operator.min.io/",
	)

	kubehelpers.CreateExternalSecret(chart, namespace, "root")

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,   // namespace
		repoName,    // repository name, same as above
		"operator",  // the chart name
		releaseName, // the release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"minio-operator.yaml",
			),
		},
		nil,
	)

	k8s.NewKubeSecret(
		chart,
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

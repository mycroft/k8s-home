package apps

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewRedmineChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "redmine"
	releaseName := namespace

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)
	kubehelpers.CreateSecretStore(chart, namespace)
	kubehelpers.CreateExternalSecret(chart, namespace, "mariadb")

	kubehelpers.CreateHelmRepository(
		chart,
		"redmine",
		"oci://registry-1.docker.io/bitnamicharts",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,   // namespace
		"redmine",   // repo name
		"redmine",   // chart name
		releaseName, // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"redmine.yaml",
			),
		},
		nil,
	)

	return chart
}

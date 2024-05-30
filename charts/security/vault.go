package security

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewVaultChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "vault"
	repoName := "hashicorp"
	releaseName := "vault"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.CreateHelmRepository(
		chart,
		repoName,
		"https://helm.releases.hashicorp.com",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repoName,
		"vault",     // chart name
		releaseName, // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName, // release name to be modified
				"vault.yaml",
			),
		},
		nil,
	)

	return chart
}

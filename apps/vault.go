package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
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

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		repoName,
		"https://helm.releases.hashicorp.com",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		repoName,
		"vault",     // chart name
		releaseName, // release name
		"0.25.0",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
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

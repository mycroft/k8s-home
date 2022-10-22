package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewVaultChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "vault"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"hashicorp",
		"https://helm.releases.hashicorp.com",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"hashicorp", // repo name
		"vault",     // chart name
		"vault",     // release name
		"0.22.0",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{},
		nil,
	)

	return chart
}

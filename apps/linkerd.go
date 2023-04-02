package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewLinkerdChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "linkerd"
	repositoryName := "linkerd"
	chartName := "linkerd-control-plane"
	releaseName := "linkerd-control-plane"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://helm.linkerd.io/stable",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName, // repo name
		"linkerd-crds", // chart name
		"linkerd-crds", // release name
		"1.4.0",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{},
		nil,
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName, // repo name
		chartName,      // chart name
		releaseName,    // release name
		"1.9.6",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"linkerd-control-plane.yaml",
			),
		},
		nil,
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName, // repo name
		"linkerd-viz",  // chart name
		"linkerd-viz",  // release name
		"30.3.6",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{},
		nil,
	)

	return chart
}

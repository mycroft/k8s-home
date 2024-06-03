package infra

import (
	"context"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewLinkerdChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	namespace := "linkerd"
	repositoryName := "linkerd"
	chartName := "linkerd-control-plane"
	releaseName := "linkerd-control-plane"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://helm.linkerd.io/stable",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName, // repo name
		"linkerd-crds", // chart name
		"linkerd-crds", // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{},
		nil,
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName, // repo name
		chartName,      // chart name
		releaseName,    // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"linkerd-control-plane.yaml",
			),
		},
		nil,
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName, // repo name
		"linkerd-viz",  // chart name
		"linkerd-viz",  // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{},
		nil,
	)

	return chart
}

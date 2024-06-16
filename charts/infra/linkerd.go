package infra

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewLinkerdChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "linkerd"
	repositoryName := "linkerd"
	chartName := "linkerd-control-plane"
	releaseName := "linkerd-control-plane"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateHelmRepository(
		chart.Cdk8sChart,
		repositoryName,
		"https://helm.linkerd.io/stable",
	)

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,
		repositoryName, // repo name
		"linkerd-crds", // chart name
		"linkerd-crds", // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{},
		nil,
	)

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,
		repositoryName, // repo name
		chartName,      // chart name
		releaseName,    // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"linkerd-control-plane.yaml",
			),
		},
		nil,
	)

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
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

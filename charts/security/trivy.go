package security

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewTrivyChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "trivy"
	namespace := "trivy-system"

	repoName := "aqua"
	releaseName := "trivy-operator"

	chart := builder.NewChart(appName)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repoName,
		"https://aquasecurity.github.io/helm-charts/",
	)

	values := map[string]*string{
		"trivy.ignoreUnfixed": jsii.String("true"),
	}

	chart.CreateHelmRelease(
		namespace,
		repoName,
		"trivy-operator",
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
		kubehelpers.WithHelmValues(values),
	)

	return chart
}

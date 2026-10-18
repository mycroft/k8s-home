package security

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewTrivyChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "trivy"
	namespace := "trivy-system"

	// Not "aqua": harbor-scanner-trivy already owns flux-system/aqua, pointing
	// at helm.aquasec.com, which does not carry trivy-operator.
	repoName := "aquasecurity"
	releaseName := "trivy-operator"

	chart := builder.NewChart(appName)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repoName,
		"https://aquasecurity.github.io/helm-charts/",
	)

	chart.CreateHelmRelease(
		namespace,
		repoName,
		"trivy-operator",
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
		// The cluster still carries CRDs from the 2023 install.
		kubehelpers.WithReplaceCRDs(),
	)

	return chart
}

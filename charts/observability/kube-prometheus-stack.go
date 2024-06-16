package observability

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewKubePrometheusStackChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "kube-prometheus-stack"
	namespace := "monitoring"

	repositoryName := "prometheus-community"
	chartName := "kube-prometheus-stack"
	releaseName := "prometheus"

	chart := builder.NewChart(appName)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)

	chart.CreateHelmRepository(
		repositoryName,
		"https://prometheus-community.github.io/helm-charts",
	)

	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "grafana-secret")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "grafana-oidc-client")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "grafana-authentik-oidc-client")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "alertmanager-config")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "grafana-postgres")

	chart.CreateHelmRelease(
		namespace,
		repositoryName, // repoName; must be in flux-system
		chartName,      // chart name
		releaseName,    // release name
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"kube-prometheus-stack.yaml",
			),
		},
		nil,
	)

	return chart
}

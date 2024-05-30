package observability

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewKubePrometheusStackChart(scope constructs.Construct) cdk8s.Chart {
	appName := "kube-prometheus-stack"
	namespace := "monitoring"

	repositoryName := "prometheus-community"
	chartName := "kube-prometheus-stack"
	releaseName := "prometheus"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.CreateSecretStore(chart, namespace)

	kubehelpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://prometheus-community.github.io/helm-charts",
	)

	kubehelpers.CreateExternalSecret(chart, namespace, "grafana-secret")
	kubehelpers.CreateExternalSecret(chart, namespace, "grafana-oidc-client")
	kubehelpers.CreateExternalSecret(chart, namespace, "alertmanager-config")
	kubehelpers.CreateExternalSecret(chart, namespace, "grafana-postgres")

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName, // repoName; must be in flux-system
		chartName,      // chart name
		releaseName,    // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"kube-prometheus-stack.yaml",
			),
		},
		nil,
	)

	return chart
}

package apps

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
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

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateSecretStore(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://prometheus-community.github.io/helm-charts",
	)

	k8s_helpers.CreateExternalSecret(chart, namespace, "grafana-secret")
	k8s_helpers.CreateExternalSecret(chart, namespace, "grafana-oidc-client")
	k8s_helpers.CreateExternalSecret(chart, namespace, "alertmanager-config")
	k8s_helpers.CreateExternalSecret(chart, namespace, "grafana-postgres")

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName, // repoName; must be in flux-system
		chartName,      // chart name
		releaseName,    // release name
		"48.1.1",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
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

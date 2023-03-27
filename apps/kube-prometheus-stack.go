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

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateSecretStore(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"prometheus-community",
		"https://prometheus-community.github.io/helm-charts",
	)

	k8s_helpers.CreateExternalSecret(chart, namespace, "grafana-secret")
	k8s_helpers.CreateExternalSecret(chart, namespace, "grafana-oidc-client")
	k8s_helpers.CreateExternalSecret(chart, namespace, "alertmanager-config")
	k8s_helpers.CreateExternalSecret(chart, namespace, "grafana-postgres")

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"prometheus-community",  // repoName; must be in flux-system
		"kube-prometheus-stack", // chart name
		"prometheus",            // release name
		"45.0.0",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"", // release name to be modified
				"kube-prometheus-stack.yaml",
			),
		},
		nil,
	)

	return chart
}

package observability

import (
	"log"
	"os"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
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
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "grafana-postgres-cnpg")

	// Alertmanager configuration: routing and receivers live in git, the
	// receivers' credentials in Vault. alertmanagerSpec.configSecret points at
	// the first Secret; alertmanagerSpec.secrets mounts the second, which the
	// config reads through *_file options.
	alertmanagerConfig, err := os.ReadFile("configs/alertmanager/alertmanager.yaml")
	if err != nil {
		log.Fatalf("Could not read alertmanager config: %v", err)
	}
	k8s.NewKubeSecret(
		chart.Cdk8sChart,
		jsii.String("alertmanager-routing"),
		&k8s.KubeSecretProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("alertmanager-routing"), // alertmanagerSpec.configSecret
			},
			StringData: &map[string]*string{
				"alertmanager.yaml": jsii.String(string(alertmanagerConfig)),
			},
		},
	)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "alertmanager-receivers")

	configMaps := []kubehelpers.HelmReleaseConfigMap{
		chart.CreateHelmValuesConfig(
			namespace,
			releaseName,
			"kube-prometheus-stack.yaml",
			nil,
		),
	}

	chart.CreateHelmRelease(
		namespace,
		repositoryName, // repoName; must be in flux-system
		chartName,      // chart name
		releaseName,    // release name
		kubehelpers.WithConfigMaps(configMaps),
	)

	return chart
}

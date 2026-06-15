package observability

import (
	"git.mkz.me/mycroft/k8s-home/imports/scrapeconfig_monitoringcoreoscom"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
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

	configMaps := []kubehelpers.HelmReleaseConfigMap{
		kubehelpers.CreateHelmValuesConfig(
			chart.Cdk8sChart,
			namespace,
			releaseName,
			"kube-prometheus-stack.yaml",
		),
	}

	chart.CreateHelmRelease(
		namespace,
		repositoryName, // repoName; must be in flux-system
		chartName,      // chart name
		releaseName,    // release name
		kubehelpers.WithConfigMaps(configMaps),
	)

	// Custom monitors
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "garage-monitoring")

	// Scrape the Garage backup target's metrics endpoint, which is not reachable via a
	// Kubernetes Service since it lives outside the cluster.
	scrapeconfig_monitoringcoreoscom.NewScrapeConfig(
		chart.Cdk8sChart,
		jsii.String("garage-monitoring"),
		&scrapeconfig_monitoringcoreoscom.ScrapeConfigProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Labels: &map[string]*string{
					"release": jsii.String("prometheus"),
				},
			},
			Spec: &scrapeconfig_monitoringcoreoscom.ScrapeConfigSpec{
				StaticConfigs: &[]*scrapeconfig_monitoringcoreoscom.ScrapeConfigSpecStaticConfigs{
					{
						Targets: &[]*string{
							jsii.String("moonstone.lan.mkz.me:3903"),
						},
						Labels: &map[string]*string{
							"job": jsii.String("garage-monitoring"),
						},
					},
				},
				MetricsPath: jsii.String("/metrics"),
				Scheme:      scrapeconfig_monitoringcoreoscom.ScrapeConfigSpecScheme_HTTP,
				Authorization: &scrapeconfig_monitoringcoreoscom.ScrapeConfigSpecAuthorization{
					Credentials: &scrapeconfig_monitoringcoreoscom.ScrapeConfigSpecAuthorizationCredentials{
						Name: jsii.String("garage-monitoring"),
						Key:  jsii.String("metrics_token"),
					},
				},
			},
		},
	)

	return chart
}

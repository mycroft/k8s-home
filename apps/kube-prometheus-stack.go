package apps

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	"git.mkz.me/mycroft/k8s-home/imports/acidzalando"
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
				"kube-prometheus-stack.yaml",
			),
		},
		nil,
	)

	// Spawn a PostgreSQL server for Grafana.
	acidzalando.NewPostgresql(
		chart,
		jsii.String(fmt.Sprintf("postgres-%s", namespace)),
		&acidzalando.PostgresqlProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String(fmt.Sprintf("postgres-%s", namespace)),
				Namespace: jsii.String(namespace),
			},
			Spec: &acidzalando.PostgresqlSpec{
				TeamId: jsii.String(namespace),
				Volume: &acidzalando.PostgresqlSpecVolume{
					StorageClass: jsii.String("longhorn"),
					Size:         jsii.String("512Mi"),
				},
				NumberOfInstances: jsii.Number(float64(1)),
				Databases: &map[string]*string{
					"grafana": jsii.String("grafana"),
				},
				Users: &map[string]*[]acidzalando.PostgresqlSpecUsers{
					"grafana-admin": {
						acidzalando.PostgresqlSpecUsers_SUPERUSER,
						acidzalando.PostgresqlSpecUsers_CREATEDB,
					},
					"grafana": {},
				},
				Postgresql: &acidzalando.PostgresqlSpecPostgresql{
					Version: acidzalando.PostgresqlSpecPostgresqlVersion_VALUE_15,
				},
			},
		},
	)

	return chart
}

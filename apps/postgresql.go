package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/acidzalando"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewPostgres(scope constructs.Construct) cdk8s.Chart {
	namespace := "postgres"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	// Spawn a PostgreSQL server for multiple databases.
	acidzalando.NewPostgresql(
		chart,
		jsii.String(fmt.Sprintf("%s-instance", namespace)),
		&acidzalando.PostgresqlProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String(fmt.Sprintf("%s-instance", namespace)),
				Namespace: jsii.String(namespace),
			},
			Spec: &acidzalando.PostgresqlSpec{
				TeamId: jsii.String(namespace),
				Volume: &acidzalando.PostgresqlSpecVolume{
					StorageClass: jsii.String("longhorn-crypto-global"),
					Size:         jsii.String("64Gi"),
				},
				NumberOfInstances: jsii.Number(float64(1)),
				Databases: &map[string]*string{
					"grafana":  jsii.String("grafana-admin"),
					"testaroo": jsii.String("testaroo-admin"),
				},
				Users: &map[string]*[]acidzalando.PostgresqlSpecUsers{
					"grafana-admin": {
						acidzalando.PostgresqlSpecUsers_SUPERUSER,
						acidzalando.PostgresqlSpecUsers_CREATEDB,
					},
					"grafana": {},
					"testaroo-admin": {
						acidzalando.PostgresqlSpecUsers_SUPERUSER,
						acidzalando.PostgresqlSpecUsers_CREATEDB,
					},
					"testaroo": {},
				},
				Postgresql: &acidzalando.PostgresqlSpecPostgresql{
					Version: acidzalando.PostgresqlSpecPostgresqlVersion_VALUE_15,
				},
			},
		},
	)

	return chart
}

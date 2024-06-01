package storage

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/acidzalando"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
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

	kubehelpers.NewNamespace(chart, namespace)

	databases := []string{
		"grafana",
		"testaroo",
		"wallabag",
		"freshrss",
		"paperlessngx",
		"dex",
		"wikijs",
		"temporal",
		"temporal_visibility",
		"backstage",
		"authentik",
	}

	databaseSpecs := map[string]*string{}
	databaseUsers := map[string]*[]acidzalando.PostgresqlSpecUsers{}
	for _, database := range databases {
		databaseSpecs[database] = jsii.String(fmt.Sprintf("%s-admin", database))
		databaseUsers[database] = &[]acidzalando.PostgresqlSpecUsers{}
		databaseUsers[fmt.Sprintf("%s-admin", database)] = &[]acidzalando.PostgresqlSpecUsers{
			acidzalando.PostgresqlSpecUsers_SUPERUSER,
			acidzalando.PostgresqlSpecUsers_CREATEDB,
		}
	}

	env := []interface{}{
		map[string]interface{}{
			"name":  jsii.String("ALLOW_NOSSL"),
			"value": jsii.String("1"),
		},
	}

	// Spawn a PostgreSQL server for multiple databases.
	// Don't forget that "users" do not have the right to change stuff in schemas.
	// Therefore you might want to do the following:
	// GRANT CREATE ON SCHEMA public TO PUBLIC;
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
				Env:               &env,
				NumberOfInstances: jsii.Number(float64(1)),
				Databases:         &databaseSpecs,
				Users:             &databaseUsers,
				Postgresql: &acidzalando.PostgresqlSpecPostgresql{
					Version: acidzalando.PostgresqlSpecPostgresqlVersion_VALUE_15,
				},
			},
		},
	)

	return chart
}

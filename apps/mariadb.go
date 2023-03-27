package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/mariadbmmontesio"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewMariaDBChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "mariadb"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)
	k8s_helpers.CreateExternalSecret(chart, namespace, "mariadb")
	k8s_helpers.CreateExternalSecret(chart, namespace, "mariadb-testaroo")

	mariadbmmontesio.NewMariaDb(
		chart,
		jsii.String("mariadb"),
		&mariadbmmontesio.MariaDbProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("mariadb"),
				Namespace: jsii.String(namespace),
			},
			Spec: &mariadbmmontesio.MariaDbSpec{
				RootPasswordSecretKeyRef: &mariadbmmontesio.MariaDbSpecRootPasswordSecretKeyRef{
					Name: jsii.String("mariadb"),
					Key:  jsii.String("root-password"),
				},
				Image: &mariadbmmontesio.MariaDbSpecImage{
					Repository: jsii.String("mariadb"),
					Tag:        jsii.String("10.7.4"),
					PullPolicy: jsii.String("IfNotPresent"),
				},
				Port: jsii.Number(3306),
				VolumeClaimTemplate: &mariadbmmontesio.MariaDbSpecVolumeClaimTemplate{
					Resources: &mariadbmmontesio.MariaDbSpecVolumeClaimTemplateResources{
						Requests: &map[string]mariadbmmontesio.MariaDbSpecVolumeClaimTemplateResourcesRequests{
							"storage": mariadbmmontesio.MariaDbSpecVolumeClaimTemplateResourcesRequests_FromString(jsii.String("8Gi")),
						},
					},
					StorageClassName: jsii.String("longhorn-crypto-global"),
					AccessModes: &[]*string{
						jsii.String("ReadWriteOnce"),
					},
				},
			},
		},
	)

	mariadbmmontesio.NewDatabase(
		chart,
		jsii.String("mariadb-database-mariadb"),
		&mariadbmmontesio.DatabaseProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("mariadb"),
				Namespace: jsii.String(namespace),
			},
			Spec: &mariadbmmontesio.DatabaseSpec{
				MariaDbRef: &mariadbmmontesio.DatabaseSpecMariaDbRef{
					Name: jsii.String("mariadb"),
				},
				CharacterSet: jsii.String("utf8"),
				Collate:      jsii.String("utf8_general_ci"),
			},
		},
	)

	mariadbmmontesio.NewUser(
		chart,
		jsii.String("mariadb-user-testaroo"),
		&mariadbmmontesio.UserProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("mariadb-testaroo"),
			},
			Spec: &mariadbmmontesio.UserSpec{
				MariaDbRef: &mariadbmmontesio.UserSpecMariaDbRef{
					Name: jsii.String("mariadb"),
				},
				PasswordSecretKeyRef: &mariadbmmontesio.UserSpecPasswordSecretKeyRef{
					Name: jsii.String("mariadb-testaroo"),
					Key:  jsii.String("password"),
				},
			},
		},
	)

	mariadbmmontesio.NewGrant(
		chart,
		jsii.String("mariadb-chart-testaroo"),
		&mariadbmmontesio.GrantProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("mariadb-testaroo"),
			},
			Spec: &mariadbmmontesio.GrantSpec{
				MariaDbRef: &mariadbmmontesio.GrantSpecMariaDbRef{
					Name: jsii.String("mariadb"),
				},
				Database: jsii.String("mariadb"),
				Username: jsii.String("mariadb-testaroo"),
				Table:    jsii.String("*"),
				Privileges: &[]*string{
					jsii.String("SELECT"),
					jsii.String("INSERT"),
					jsii.String("UPDATE"),
					jsii.String("CREATE"),
					jsii.String("ALTER"),
					jsii.String("DELETE"),
					jsii.String("DROP"),
					jsii.String("INDEX"),
				},
				GrantOption: jsii.Bool(true),
			},
		},
	)

	return chart
}

package storage

import (
	"fmt"

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

	mariadbInstance := "mariadb"

	mariadbmmontesio.NewMariaDb(
		chart,
		jsii.String("mariadb"),
		&mariadbmmontesio.MariaDbProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String(mariadbInstance),
				Namespace: jsii.String(namespace),
			},
			Spec: &mariadbmmontesio.MariaDbSpec{
				Env: &[]*mariadbmmontesio.MariaDbSpecEnv{
					{
						Name:  jsii.String("MARIADB_AUTO_UPGRADE"),
						Value: jsii.String("1"),
					},
				},
				RootPasswordSecretKeyRef: &mariadbmmontesio.MariaDbSpecRootPasswordSecretKeyRef{
					Name: jsii.String("mariadb"),
					Key:  jsii.String("root-password"),
				},
				Image: jsii.String("mariadb:10.11"),
				Port:  jsii.Number(3306),
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

	databases := []string{
		"bookstack",
		"mariadb",
	}

	users := map[string][]string{
		"bookstack": []string{
			"bookstack",
		},
		"mariadb": []string{
			"mariadb",
		},
	}

	for _, database := range databases {
		mariadbmmontesio.NewDatabase(
			chart,
			jsii.String(fmt.Sprintf("%s-database", database)),
			&mariadbmmontesio.DatabaseProps{
				Metadata: &cdk8s.ApiObjectMetadata{
					Namespace: jsii.String(namespace),
					Name:      jsii.String(database),
				},
				Spec: &mariadbmmontesio.DatabaseSpec{
					MariaDbRef: &mariadbmmontesio.DatabaseSpecMariaDbRef{
						Name: jsii.String(mariadbInstance),
					},
					CharacterSet: jsii.String("utf8"),
					Collate:      jsii.String("utf8_general_ci"),
				},
			},
		)

		// users names are the same than databases (for now).
		for _, user := range users[database] {
			k8s_helpers.CreateExternalSecret(chart, namespace, fmt.Sprintf("user-%s", user))
			mariadbmmontesio.NewUser(
				chart,
				jsii.String(fmt.Sprintf("%s-user-%s", database, user)),
				&mariadbmmontesio.UserProps{
					Metadata: &cdk8s.ApiObjectMetadata{
						Namespace: jsii.String(namespace),
						Name:      jsii.String(user),
					},
					Spec: &mariadbmmontesio.UserSpec{
						MariaDbRef: &mariadbmmontesio.UserSpecMariaDbRef{
							Name: jsii.String(mariadbInstance),
						},
						PasswordSecretKeyRef: &mariadbmmontesio.UserSpecPasswordSecretKeyRef{
							Name: jsii.String(fmt.Sprintf("user-%s", user)),
							Key:  jsii.String("password"),
						},
					},
				},
			)

			mariadbmmontesio.NewGrant(
				chart,
				jsii.String(fmt.Sprintf("%s-grant-%s", database, user)),
				&mariadbmmontesio.GrantProps{
					Metadata: &cdk8s.ApiObjectMetadata{
						Namespace: jsii.String(namespace),
						Name:      jsii.String(user),
					},
					Spec: &mariadbmmontesio.GrantSpec{
						MariaDbRef: &mariadbmmontesio.GrantSpecMariaDbRef{
							Name: jsii.String(mariadbInstance),
						},
						Database: jsii.String(database),
						Username: jsii.String(user),
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
		}
	}

	return chart
}

package storage

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8smariadbcom"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
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

	kubehelpers.NewNamespace(chart, namespace)
	kubehelpers.CreateSecretStore(chart, namespace)
	kubehelpers.CreateExternalSecret(chart, namespace, "mariadb")

	mariadbInstance := "mariadb"

	k8smariadbcom.NewMariaDb(
		chart,
		jsii.String("mariadb"),
		&k8smariadbcom.MariaDbProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String(mariadbInstance),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8smariadbcom.MariaDbSpec{
				Env: &[]*k8smariadbcom.MariaDbSpecEnv{
					{
						Name:  jsii.String("MARIADB_AUTO_UPGRADE"),
						Value: jsii.String("1"),
					},
				},
				RootPasswordSecretKeyRef: &k8smariadbcom.MariaDbSpecRootPasswordSecretKeyRef{
					Name: jsii.String("mariadb"),
					Key:  jsii.String("root-password"),
				},
				Image: jsii.String("mariadb:11.3.2"),
				Port:  jsii.Number(3306),
				Storage: &k8smariadbcom.MariaDbSpecStorage{
					StorageClassName: jsii.String("longhorn-crypto-global"),
					Size:             k8smariadbcom.MariaDbSpecStorageSize_FromString(jsii.String("8Gi")),
				},
			},
		},
	)

	databases := []string{
		"bookstack",
		"mariadb",
		"redmine",
	}

	users := map[string][]string{
		"bookstack": []string{
			"bookstack",
		},
		"mariadb": []string{
			"mariadb",
		},
		"redmine": []string{
			"redmine",
		},
	}

	for _, database := range databases {
		k8smariadbcom.NewDatabase(
			chart,
			jsii.String(fmt.Sprintf("%s-database", database)),
			&k8smariadbcom.DatabaseProps{
				Metadata: &cdk8s.ApiObjectMetadata{
					Namespace: jsii.String(namespace),
					Name:      jsii.String(database),
				},
				Spec: &k8smariadbcom.DatabaseSpec{
					MariaDbRef: &k8smariadbcom.DatabaseSpecMariaDbRef{
						Name: jsii.String(mariadbInstance),
					},
					CharacterSet: jsii.String("utf8"),
					Collate:      jsii.String("utf8_general_ci"),
				},
			},
		)

		// users names are the same than databases (for now).
		for _, user := range users[database] {
			kubehelpers.CreateExternalSecret(chart, namespace, fmt.Sprintf("user-%s", user))
			k8smariadbcom.NewUser(
				chart,
				jsii.String(fmt.Sprintf("%s-user-%s", database, user)),
				&k8smariadbcom.UserProps{
					Metadata: &cdk8s.ApiObjectMetadata{
						Namespace: jsii.String(namespace),
						Name:      jsii.String(user),
					},
					Spec: &k8smariadbcom.UserSpec{
						MariaDbRef: &k8smariadbcom.UserSpecMariaDbRef{
							Name: jsii.String(mariadbInstance),
						},
						PasswordSecretKeyRef: &k8smariadbcom.UserSpecPasswordSecretKeyRef{
							Name: jsii.String(fmt.Sprintf("user-%s", user)),
							Key:  jsii.String("password"),
						},
					},
				},
			)

			k8smariadbcom.NewGrant(
				chart,
				jsii.String(fmt.Sprintf("%s-grant-%s", database, user)),
				&k8smariadbcom.GrantProps{
					Metadata: &cdk8s.ApiObjectMetadata{
						Namespace: jsii.String(namespace),
						Name:      jsii.String(user),
					},
					Spec: &k8smariadbcom.GrantSpec{
						MariaDbRef: &k8smariadbcom.GrantSpecMariaDbRef{
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

package storage

import (
	"git.mkz.me/mycroft/k8s-home/imports/cnpg_cluster_postgresqlcnpgio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// databases hosted on the CloudNativePG cluster. Adding one is a single entry:
// the database, its owning role, a generated password and the Vault credentials
// the application reads are all derived from it. Namespace only needs to be set
// when the application namespace differs from the database name.
var databases = []kubehelpers.CNPGDatabase{
	{Name: "zipline", VaultEntry: "postgresql-cnpg"},
	{Name: "memos", VaultEntry: "postgres-cnpg"},
	{Name: "vikunja", VaultEntry: "postgresql-cnpg"},
	{Name: "wikijs", VaultEntry: "postgresql-cnpg"},
	{Name: "wallabag", VaultEntry: "postgresql-cnpg"},
	{Name: "privatebin", VaultEntry: "postgres-cnpg"},
	{Name: "n8n", VaultEntry: "postgresql-cnpg", PasswordAliases: []string{"postgres-password"}},
	{Name: "grafana", Namespace: "monitoring", VaultEntry: "grafana-postgres-cnpg"},
	{Name: "dex", Namespace: "dex-idp", VaultEntry: "postgresql-cnpg"},
	{Name: "temporal", VaultEntry: "postgresql-cnpg"},
	{Name: "temporal_visibility", Namespace: "temporal", VaultEntry: "postgresql-visibility-cnpg", Extensions: []string{"btree_gin"}},
	{Name: "freshrss", VaultEntry: "postgresql-cnpg"},
	{Name: "outline", VaultEntry: "postgres-cnpg", Extensions: []string{"pg_trgm", "unaccent", "uuid-ossp"}},
	{Name: "authentik", VaultEntry: "postgresql-cnpg"},
	{Name: "paperlessngx", Namespace: "paperless-ngx", VaultEntry: "postgresql-cnpg"},
}

func NewCNPGCluster(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "cnpg"
	clusterName := "postgres"
	clusterHost := clusterName + "-rw." + namespace

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	// The PushSecret publishing credentials to Vault runs from this namespace.
	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)

	managedRoles := make([]*cnpg_cluster_postgresqlcnpgio.ClusterSpecManagedRoles, 0, len(databases))
	for _, db := range databases {
		managedRoles = append(managedRoles, kubehelpers.NewCNPGManagedRole(db))
	}

	cluster := cnpg_cluster_postgresqlcnpgio.NewCluster(
		chart.Cdk8sChart,
		jsii.String("cnpg-cluster"),
		&cnpg_cluster_postgresqlcnpgio.ClusterProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String(clusterName),
			},
			Spec: &cnpg_cluster_postgresqlcnpgio.ClusterSpec{
				Bootstrap: &cnpg_cluster_postgresqlcnpgio.ClusterSpecBootstrap{
					Initdb: &cnpg_cluster_postgresqlcnpgio.ClusterSpecBootstrapInitdb{
						Owner:    jsii.String("postgres"),
						Database: jsii.String("postgres"),
					},
				},
				Managed: &cnpg_cluster_postgresqlcnpgio.ClusterSpecManaged{
					Roles: &managedRoles,
				},
				EnableSuperuserAccess: jsii.Bool(true),
				Instances:             jsii.Number(2),
				Storage: &cnpg_cluster_postgresqlcnpgio.ClusterSpecStorage{
					StorageClass: jsii.String("longhorn-crypto-global"),
					Size:         jsii.String("64Gi"),
				},
			},
		},
	)

	for _, db := range databases {
		kubehelpers.NewCNPGDatabase(chart.Cdk8sChart, kubehelpers.CNPGDatabaseConfig{
			Namespace:   namespace,
			ClusterName: cluster.Name(),
			Host:        clusterHost,
			Database:    db,
		})
	}

	return chart
}

package storage

import (
	"git.mkz.me/mycroft/k8s-home/imports/cnpg_cluster_postgresqlcnpgio"
	database "git.mkz.me/mycroft/k8s-home/imports/cnpg_database_postgresqlcnpgio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewCNPGCluster(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "cnpg"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	managedRoles := []*cnpg_cluster_postgresqlcnpgio.ClusterSpecManagedRoles{
		{
			Name:            jsii.String("zipline"),
			Ensure:          cnpg_cluster_postgresqlcnpgio.ClusterSpecManagedRolesEnsure_PRESENT,
			Createdb:        jsii.Bool(false),
			Createrole:      jsii.Bool(false),
			DisablePassword: jsii.Bool(false),
			Login:           jsii.Bool(true),
			Superuser:       jsii.Bool(false),
		},
	}

	cluster := cnpg_cluster_postgresqlcnpgio.NewCluster(
		chart.Cdk8sChart,
		jsii.String("cnpg-cluster"),
		&cnpg_cluster_postgresqlcnpgio.ClusterProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("postgres"),
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

	database.NewDatabase(
		chart.Cdk8sChart,
		jsii.String("zipline"),
		&database.DatabaseProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.Sprintf("postgres-zipline"),
			},
			Spec: &database.DatabaseSpec{
				Cluster: &database.DatabaseSpecCluster{
					Name: cluster.Name(),
				},
				Name:  jsii.String("zipline"),
				Owner: jsii.String("zipline"),
			},
		},
	)

	return chart
}

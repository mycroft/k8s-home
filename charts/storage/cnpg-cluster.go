package storage

import (
	"git.mkz.me/mycroft/k8s-home/imports/cnpg_cluster_postgresqlcnpgio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewCNPGCluster(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "cnpg"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	cnpg_cluster_postgresqlcnpgio.NewCluster(
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
				EnableSuperuserAccess: jsii.Bool(true),
				Instances:             jsii.Number(2),
				Storage: &cnpg_cluster_postgresqlcnpgio.ClusterSpecStorage{
					StorageClass: jsii.String("longhorn-crypto-global"),
					Size:         jsii.String("64Gi"),
				},
			},
		},
	)

	return chart
}

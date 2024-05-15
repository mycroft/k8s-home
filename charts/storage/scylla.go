package storage

import (
	"git.mkz.me/mycroft/k8s-home/imports/scyllascylladbcom"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewScyllaChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "scylla"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	// See https://operator.docs.scylladb.com/stable/scylla_cluster_crd.html
	// Sample https://github.com/scylladb/scylla-operator/blob/master/examples/generic/cluster.yaml
	scyllascylladbcom.NewScyllaCluster(
		chart,
		jsii.String("scylla-cluster"),
		&scyllascylladbcom.ScyllaClusterProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("scylla-cluster"),
				Namespace: jsii.String("scylla"),
			},
			Spec: &scyllascylladbcom.ScyllaClusterSpec{
				Version:       jsii.String("5.4.3"),
				AgentVersion:  jsii.String("3.2.6"),
				Repository:    jsii.String("scylladb/scylla"),
				DeveloperMode: jsii.Bool(true), // using DeveloperMode to bypass FS checks
				Datacenter: &scyllascylladbcom.ScyllaClusterSpecDatacenter{
					Name: jsii.String("eu-west-1"),
					Racks: &[]*scyllascylladbcom.ScyllaClusterSpecDatacenterRacks{
						{
							Name:    jsii.String("eu-west-1a"),
							Members: jsii.Number(3),
							Storage: &scyllascylladbcom.ScyllaClusterSpecDatacenterRacksStorage{
								Capacity:         jsii.String("32Gi"),
								StorageClassName: jsii.String("longhorn-crypto-global"),
							},
							Resources: &scyllascylladbcom.ScyllaClusterSpecDatacenterRacksResources{
								Requests: &map[string]scyllascylladbcom.ScyllaClusterSpecDatacenterRacksResourcesRequests{
									"cpu":    scyllascylladbcom.ScyllaClusterSpecDatacenterRacksResourcesRequests_FromString(jsii.String("0.25")),
									"memory": scyllascylladbcom.ScyllaClusterSpecDatacenterRacksResourcesRequests_FromString(jsii.String("2Gi")),
								},
								Limits: &map[string]scyllascylladbcom.ScyllaClusterSpecDatacenterRacksResourcesLimits{
									"cpu":    scyllascylladbcom.ScyllaClusterSpecDatacenterRacksResourcesLimits_FromString(jsii.String("1")),
									"memory": scyllascylladbcom.ScyllaClusterSpecDatacenterRacksResourcesLimits_FromString(jsii.String("4Gi")),
								},
							},
						},
					},
				},
			},
		},
	)

	return chart
}

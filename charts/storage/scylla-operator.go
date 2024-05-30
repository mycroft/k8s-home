package storage

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// CRDs are:
//
// https://raw.githubusercontent.com/scylladb/scylla-operator/master/pkg/api/scylla/v1alpha1/scylla.scylladb.com_scylladbmonitorings.yaml
// https://raw.githubusercontent.com/scylladb/scylla-operator/master/pkg/api/scylla/v1alpha1/scylla.scylladb.com_nodeconfigs.yaml
// https://raw.githubusercontent.com/scylladb/scylla-operator/master/pkg/api/scylla/v1alpha1/scylla.scylladb.com_scyllaoperatorconfigs.yaml
// https://raw.githubusercontent.com/scylladb/scylla-operator/master/pkg/api/scylla/v1/scylla.scylladb.com_scyllaclusters.yaml
//
// TODO: Automatize CRDs checks/updates

func NewScyllaOperatorChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "scylla-operator"

	repoName := "scylla"
	releaseName := "scylla-operator"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.CreateHelmRepository(
		chart,
		repoName,
		"https://scylla-operator-charts.storage.googleapis.com/stable",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repoName,
		"scylla-operator",
		releaseName,
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"scylla-operator.yaml",
			),
		},
		nil,
	)

	return chart
}

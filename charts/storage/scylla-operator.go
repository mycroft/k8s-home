package storage

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
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

func NewScyllaOperatorChart(builder *kubehelpers.Builder) cdk8s.Chart {
	namespace := "scylla-operator"

	repoName := "scylla"
	releaseName := "scylla-operator"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateHelmRepository(
		chart.Cdk8sChart,
		repoName,
		"https://scylla-operator-charts.storage.googleapis.com/stable",
	)

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,
		repoName,
		"scylla-operator",
		releaseName,
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"scylla-operator.yaml",
			),
		},
		nil,
	)

	return chart.Cdk8sChart
}

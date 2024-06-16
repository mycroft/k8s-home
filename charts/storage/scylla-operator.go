package storage

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

// CRDs are:
//
// https://raw.githubusercontent.com/scylladb/scylla-operator/master/pkg/api/scylla/v1alpha1/scylla.scylladb.com_scylladbmonitorings.yaml
// https://raw.githubusercontent.com/scylladb/scylla-operator/master/pkg/api/scylla/v1alpha1/scylla.scylladb.com_nodeconfigs.yaml
// https://raw.githubusercontent.com/scylladb/scylla-operator/master/pkg/api/scylla/v1alpha1/scylla.scylladb.com_scyllaoperatorconfigs.yaml
// https://raw.githubusercontent.com/scylladb/scylla-operator/master/pkg/api/scylla/v1/scylla.scylladb.com_scyllaclusters.yaml
//
// TODO: Automatize CRDs checks/updates

func NewScyllaOperatorChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "scylla-operator"

	repoName := "scylla"
	releaseName := "scylla-operator"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repoName,
		"https://scylla-operator-charts.storage.googleapis.com/stable",
	)

	chart.CreateHelmRelease(
		namespace,
		repoName,
		"scylla-operator",
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	return chart
}

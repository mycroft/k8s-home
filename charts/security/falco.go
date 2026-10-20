package security

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

// NewFalcoChart deploys Falco (https://falco.org), which watches syscalls on
// every node through its modern eBPF probe and reports rule matches as JSON
// on stdout. Alloy ships them to Loki with the other pod logs, and Falco's
// metrics count the matches per rule. Nothing alerts on them yet: the default
// rules are first tuned against what the cluster's own workloads normally do.
// Values, and the rule exceptions as they are written, are in
// configs/falco.yaml.
func NewFalcoChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "falco"
	repositoryName := "falcosecurity"
	chartName := "falco"
	releaseName := "falco"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repositoryName,
		"https://falcosecurity.github.io/charts",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	return chart
}

package observability

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewGrafanaHelmRepositoryChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	repositoryName := "grafana"

	chart := builder.NewChart("grafana-helm-repository")

	chart.CreateHelmRepository(
		repositoryName,
		"https://grafana.github.io/helm-charts",
	)

	return chart
}

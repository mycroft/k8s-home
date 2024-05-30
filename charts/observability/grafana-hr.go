package observability

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewGrafanaHelmRepositoryChart(scope constructs.Construct) cdk8s.Chart {
	repositoryName := "grafana"

	chart := cdk8s.NewChart(
		scope,
		jsii.String("grafana-helm-repository"),
		&cdk8s.ChartProps{},
	)

	kubehelpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://grafana.github.io/helm-charts",
	)

	return chart
}

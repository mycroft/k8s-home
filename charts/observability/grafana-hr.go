package observability

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
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

	k8s_helpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://grafana.github.io/helm-charts",
	)

	return chart
}

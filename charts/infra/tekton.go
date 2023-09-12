package infra

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewTektonChart(scope constructs.Construct) cdk8s.Chart {
	chart := cdk8s.NewChart(
		scope,
		jsii.String("tekton"),
		&cdk8s.ChartProps{},
	)

	cdk8s.NewInclude(chart, jsii.String("tekton-pipeline"), &cdk8s.IncludeProps{
		Url: jsii.String("charts/static/tekton/tekton-pipeline-release.yaml"),
	})

	cdk8s.NewInclude(chart, jsii.String("tekton-triggers-release"), &cdk8s.IncludeProps{
		Url: jsii.String("charts/static/tekton/tekton-triggers-release.yaml"),
	})

	cdk8s.NewInclude(chart, jsii.String("tekton-triggers-interceptors"), &cdk8s.IncludeProps{
		Url: jsii.String("charts/static/tekton/tekton-triggers-interceptors.yaml"),
	})

	cdk8s.NewInclude(chart, jsii.String("tekton-dashboard"), &cdk8s.IncludeProps{
		Url: jsii.String("charts/static/tekton/tekton-dashboard-release.yaml"),
	})

	return chart
}

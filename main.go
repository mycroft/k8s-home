package main

import (
	"example.com/k8s-home/imports/k8s"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type MyChartProps struct {
	cdk8s.ChartProps
}

func NewMyChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	k8s.NewKubeNamespace(
		chart,
		jsii.String("ns"),
		&k8s.KubeNamespaceProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String("testaroo"),
			},
		},
	)

	// define resources here

	return chart
}

func main() {
	app := cdk8s.NewApp(nil)
	NewMyChart(app, "k8s-home", nil)
	app.Synth()
}

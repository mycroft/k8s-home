package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewHelloKubernetesChart(scope constructs.Construct) cdk8s.Chart {
	appName := "hello-kubernetes"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	k8s.NewKubeNamespace(
		chart,
		jsii.String("ns"),
		&k8s.KubeNamespaceProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String(appName),
			},
		},
	)

	return chart
}

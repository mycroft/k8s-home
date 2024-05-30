package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewITToolsChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "it-tools"
	appIngress := "it-tools.services.mkz.me"
	appName := "it-tools"
	appPort := 80

	image := kubehelpers.RegisterDockerImage("corentinth/it-tools")

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	appLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("it-tools"),
	}

	kubehelpers.NewAppDeployment(
		chart,
		namespace,
		appName,
		image,
		appLabels,
		[]*k8s.EnvVar{},
		[]string{},
		[]kubehelpers.ConfigMapMount{},
	)

	kubehelpers.NewAppIngress(
		chart,
		appLabels,
		appName,
		appPort,
		appIngress,
		map[string]string{},
	)

	return chart
}

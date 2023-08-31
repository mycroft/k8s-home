package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewITToolsChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "it-tools"
	appIngress := "it-tools.services.mkz.me"
	appName := "it-tools"
	appPort := 80

	image := k8s_helpers.RegisterDockerImage("corentinth/it-tools:2023.8.21-6f93cba")

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	appLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("it-tools"),
	}

	k8s_helpers.NewAppDeployment(
		chart,
		namespace,
		appName,
		image,
		appLabels,
		[]*k8s.EnvVar{},
		[]string{},
		[]k8s_helpers.ConfigMapMount{},
	)

	k8s_helpers.NewAppIngress(
		chart,
		appLabels,
		appName,
		appPort,
		appIngress,
		map[string]string{},
	)

	return chart
}

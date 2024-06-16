package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewITToolsChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "it-tools"
	appIngress := "it-tools.services.mkz.me"
	appName := "it-tools"
	appPort := 80

	image := kubehelpers.RegisterDockerImage("corentinth/it-tools")

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	appLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("it-tools"),
	}

	kubehelpers.NewAppDeployment(
		chart.Cdk8sChart,
		namespace,
		appName,
		image,
		appLabels,
		[]*k8s.EnvVar{},
		[]string{},
		[]kubehelpers.ConfigMapMount{},
	)

	kubehelpers.NewAppIngress(
		builder.Context,
		chart.Cdk8sChart,
		appLabels,
		appName,
		appPort,
		appIngress,
		"",
		map[string]string{},
	)

	return chart
}

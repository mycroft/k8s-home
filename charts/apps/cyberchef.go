package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewCyberchefChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "cyberchef"
	namespace := appName
	appImage := kubehelpers.RegisterDockerImage("ghcr.io/gchq/cyberchef")
	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}
	env := []*k8s.EnvVar{}

	appPort := 80
	appIngress := "cyberchef.services.mkz.me"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.NewAppDeployment(
		chart.Cdk8sChart,
		namespace,
		appName,
		appImage,
		labels,
		env,
		[]string{},
		[]kubehelpers.ConfigMapMount{},
	)

	kubehelpers.NewAppIngress(
		builder.Context,
		chart.Cdk8sChart,
		labels,
		appName,
		appPort,
		appIngress,
		"",
		map[string]string{},
	)

	return chart
}

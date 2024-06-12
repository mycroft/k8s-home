package apps

import (
	"context"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewCyberchefChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	appName := "cyberchef"
	namespace := appName
	appImage := kubehelpers.RegisterDockerImage("ghcr.io/gchq/cyberchef")
	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}
	env := []*k8s.EnvVar{}

	appPort := 80
	appIngress := "cyberchef.services.mkz.me"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.NewAppDeployment(
		chart,
		namespace,
		appName,
		appImage,
		labels,
		env,
		[]string{},
		[]kubehelpers.ConfigMapMount{},
	)

	kubehelpers.NewAppIngress(
		ctx,
		chart,
		labels,
		appName,
		appPort,
		appIngress,
		"",
		map[string]string{},
	)

	return chart
}

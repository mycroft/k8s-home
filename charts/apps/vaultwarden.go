package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewVaultWardenChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "vaultwarden"
	appName := namespace
	appImage := kubehelpers.RegisterDockerImage("vaultwarden/server")
	appPort := 80
	appIngress := "vaultwarden.services.mkz.me"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.NewStatefulSet(
		chart,
		namespace,
		appName,
		appImage,
		appPort,
		labels,
		[]*k8s.EnvVar{},
		[]string{},
		[]kubehelpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/data",
				StorageSize: "10Gi",
			},
		},
	)

	kubehelpers.NewAppIngress(
		chart,
		labels,
		appName,
		appPort,
		appIngress,
		map[string]string{},
	)

	return chart
}

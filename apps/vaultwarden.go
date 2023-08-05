package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewVaultWardenChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "vaultwarden"
	appName := namespace
	appImage := k8s_helpers.RegisterDockerImage("vaultwarden/server:1.29.1")
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

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.NewStatefulSet(
		chart,
		namespace,
		appName,
		appImage,
		appPort,
		labels,
		[]*k8s.EnvVar{},
		[]string{},
		[]k8s_helpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/data",
				StorageSize: "10Gi",
			},
		},
	)

	k8s_helpers.NewAppIngress(
		chart,
		labels,
		appName,
		appPort,
		appIngress,
		map[string]string{},
	)

	return chart
}

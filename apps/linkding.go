package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

const (
	linkdingImage = "sissbruecker/linkding"
)

func NewLinkdingChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "linkding"
	appName := namespace
	appPort := 9090
	appIngress := "links.services.mkz.me"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	k8s_helpers.NewStatefulSet(
		chart,
		namespace,
		appName,
		linkdingImage,
		appPort,
		labels,
		[]*k8s.EnvVar{},
		[]string{},
		[]k8s_helpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/etc/linkding/data",
				StorageSize: "1Gi",
			},
		},
		true, // use legacy naming
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
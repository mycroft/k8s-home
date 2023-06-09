package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewSnippetBoxChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "snippetbox"
	appIngress := "snippetbox.services.mkz.me"
	appName := "snippetbox"
	appPort := 5000
	image := k8s_helpers.RegisterDockerImage("pawelmalak/snippet-box:latest")

	labels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("snippetbox"),
	}

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.NewStatefulSet(
		chart,
		namespace,
		appName,
		image,
		appPort,
		labels,
		[]*k8s.EnvVar{},
		[]string{},
		[]k8s_helpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/app/data",
				StorageSize: "1Gi",
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

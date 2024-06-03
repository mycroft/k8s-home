package apps

import (
	"context"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewLinkdingChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	namespace := "linkding"
	appName := namespace
	appPort := 9090
	appIngress := "links.services.mkz.me"
	linkdingImage := kubehelpers.RegisterDockerImage("sissbruecker/linkding")

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)
	kubehelpers.CreateSecretStore(chart, namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	_, svcName := kubehelpers.NewStatefulSet(
		chart,
		namespace,
		appName,
		linkdingImage,
		appPort,
		labels,
		[]*k8s.EnvVar{},
		[]string{},
		[]kubehelpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/etc/linkding/data",
				StorageSize: "1Gi",
			},
		},
	)

	kubehelpers.NewAppIngress(
		ctx,
		chart,
		labels,
		appName,
		appPort,
		appIngress,
		svcName,
		map[string]string{},
	)

	return chart
}

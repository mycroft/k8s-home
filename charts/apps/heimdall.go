package apps

import (
	"context"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewHeimdallChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	namespace := "heimdall"
	appName := "heimdall"
	appImage := kubehelpers.RegisterDockerImage("linuxserver/heimdall")
	appPort := 80
	appIngress := "heimdall.services.mkz.me"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	env := []*k8s.EnvVar{
		{Name: jsii.String("PUID"), Value: jsii.String("1000")},
		{Name: jsii.String("PGID"), Value: jsii.String("1000")},
		{Name: jsii.String("TZ"), Value: jsii.String("Etc/UTC")},
	}

	kubehelpers.NewStatefulSet(
		chart,
		namespace,
		appName,
		appImage,
		appPort,
		labels,
		env,
		[]string{},
		[]kubehelpers.StatefulSetVolume{
			{
				Name:        "config",
				MountPath:   "/config",
				StorageSize: "1Gi",
			},
			{
				Name:        "extensions",
				MountPath:   "/var/www/FreshRSS/extensions",
				StorageSize: "1Gi",
			},
		},
	)

	kubehelpers.NewAppIngress(
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

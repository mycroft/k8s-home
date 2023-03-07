package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewFreshRSS(scope constructs.Construct) cdk8s.Chart {
	namespace := "freshrss"
	appName := namespace
	appImage := "freshrss/freshrss:1.21.0"
	appPort := 80
	appIngress := "freshrss.services.mkz.me"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	env := []*k8s.EnvVar{
		{Name: jsii.String("PUID"), Value: jsii.String("1000")},
		{Name: jsii.String("PGID"), Value: jsii.String("1000")},
		{Name: jsii.String("TZ"), Value: jsii.String("Etc/UTC")},
	}

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.NewStatefulSet(
		chart,
		namespace,
		appName,
		appImage,
		appPort,
		labels,
		env,
		[]k8s_helpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/var/www/FreshRSS/data",
				StorageSize: "1Gi",
			},
			{
				Name:        "extensions",
				MountPath:   "/var/www/FreshRSS/extensions",
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
	)

	return chart
}

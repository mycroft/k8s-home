package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewHeimdallChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "heimdall"
	appName := "heimdall"
	appImage := builder.RegisterContainerImage("linuxserver/heimdall")
	appPort := 80
	appIngress := "heimdall.services.mkz.me"

	chart := builder.NewChart(appName)
	chart.NewNamespace(namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	env := []*k8s.EnvVar{
		{Name: jsii.String("PUID"), Value: jsii.String("1000")},
		{Name: jsii.String("PGID"), Value: jsii.String("1000")},
		{Name: jsii.String("TZ"), Value: jsii.String("Etc/UTC")},
	}

	kubehelpers.NewStatefulSet(
		chart.Cdk8sChart,
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

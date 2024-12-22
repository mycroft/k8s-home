package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewLinkdingChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "linkding"
	appName := namespace
	appPort := 9090
	appIngress := "links.services.mkz.me"
	linkdingImage := builder.RegisterContainerImage("sissbruecker/linkding")

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	_, svcName := kubehelpers.NewStatefulSet(
		chart.Cdk8sChart,
		namespace,
		appName,
		linkdingImage,
		appPort,
		labels,
		[]*k8s.EnvVar{},
		[]string{},
		[]kubehelpers.ConfigMapMount{},
		[]kubehelpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/etc/linkding/data",
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
		svcName,
		map[string]string{},
	)

	return chart
}

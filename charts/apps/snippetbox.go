package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"

	"github.com/aws/jsii-runtime-go"
)

func NewSnippetBoxChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "snippetbox"
	appIngress := "snippetbox.services.mkz.me"
	appName := "snippetbox"
	appPort := 5000
	image := builder.RegisterContainerImage("pawelmalak/snippet-box")

	labels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("snippetbox"),
	}

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	_, svcName := kubehelpers.NewStatefulSet(
		chart.Cdk8sChart,
		namespace,
		appName,
		image,
		appPort,
		labels,
		[]*k8s.EnvVar{},
		[]string{},
		[]kubehelpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/app/data",
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

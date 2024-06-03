package apps

import (
	"context"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewSnippetBoxChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	namespace := "snippetbox"
	appIngress := "snippetbox.services.mkz.me"
	appName := "snippetbox"
	appPort := 5000
	image := kubehelpers.RegisterDockerImage("pawelmalak/snippet-box")

	labels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("snippetbox"),
	}

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	kubehelpers.NewStatefulSet(
		chart,
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
		ctx,
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

package apps

import (
	"os"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewKubernetesDashboardChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "kubernetes-dashboard"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"kubernetes-dashboard",
		"https://kubernetes.github.io/dashboard/",
	)

	contents, err := os.ReadFile("configs/kubernetes-dashboard.yaml")
	if err != nil {
		panic(err)
	}

	cm := k8s.NewKubeConfigMap(
		chart,
		jsii.String("cm"),
		&k8s.KubeConfigMapProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
			},
			Data: &map[string]*string{
				"kubernetes-dashboard.yaml": jsii.String(string(contents)),
			},
		},
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"kubernetes-dashboard",
		"kubernetes-dashboard", // chart name
		"kubernetes-dashboard", // release name
		"5.11.0",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			{
				Name:    *cm.Name(),
				KeyName: "kubernetes_dashboard.yaml",
			},
		},
		map[string]*string{
			"configMapHash": jsii.String(k8s_helpers.ComputeConfigMapHash(cm)),
		},
	)

	return chart
}

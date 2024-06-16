package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

const (
	urlsImage = "git.mkz.me/mycroft/urls:latest"
)

func NewUrlsChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "urls"
	appName := namespace
	ingressHost := fmt.Sprintf("%s.services.mkz.me", appName)

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	k8s.NewKubeDeployment(
		chart.Cdk8sChart,
		jsii.String("deploy"),
		&k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(appName),
			},
			Spec: &k8s.DeploymentSpec{
				Selector: &k8s.LabelSelector{
					MatchLabels: &labels,
				},
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels: &labels,
					},
					Spec: &k8s.PodSpec{
						Containers: &[]*k8s.Container{
							{
								ImagePullPolicy: jsii.String("Always"),
								Name:            jsii.String(appName),
								Image:           jsii.String(urlsImage),
							},
						},
					},
				},
			},
		},
	)

	kubehelpers.NewAppIngress(
		builder.Context,
		chart.Cdk8sChart,
		labels,
		appName,
		3000,
		ingressHost,
		"",
		map[string]string{},
	)

	return chart
}

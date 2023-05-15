package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

const (
	heyImage = "git.mkz.me/mycroft/hey:0.0.5"
)

func NewHeyChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "hey"
	appName := namespace
	appPort := 3000
	ingressHost := fmt.Sprintf("%s.services.mkz.me", appName)

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	k8s.NewKubeDeployment(
		chart,
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
								Image:           jsii.String(heyImage),
							},
						},
					},
				},
			},
		},
	)

	k8s_helpers.NewAppIngress(
		chart,
		labels,
		appName,
		appPort,
		ingressHost,
		map[string]string{},
	)

	return chart
}

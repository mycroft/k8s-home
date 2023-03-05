package k8s_helpers

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewAppDeployment(
	chart cdk8s.Chart,
	namespace, appName, appImage string,
	labels map[string]*string,
	env []*k8s.EnvVar,
) {
	k8s.NewKubeDeployment(
		chart,
		jsii.String("deploy"),
		&k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
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
								Name:  jsii.String(appName),
								Image: jsii.String(appImage),
								Env:   &env,
							},
						},
					},
				},
			},
		},
	)
}

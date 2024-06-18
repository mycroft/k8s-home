package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewGossipChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	name := "gossip"

	namespace := name
	appName := name
	imageName := "registry.mkz.me/mycroft/gossip:latest"

	chart := builder.NewChart(name)
	chart.NewNamespace(namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	env := []*k8s.EnvVar{
		{
			Name: jsii.String("NATS_SERVER"),
			// there is something wierd with the headless service.
			Value: jsii.String("nats.nats:4222"),
		},
		{
			Name:  jsii.String("DELAY"),
			Value: jsii.String("10"),
		},
	}

	k8s.NewKubeDeployment(
		chart.Cdk8sChart,
		jsii.String("deploy-server"),
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
								Env:             &env,
								Name:            jsii.String(appName),
								Image:           jsii.String(imageName),
								ImagePullPolicy: jsii.String("Always"),
								Command: &[]*string{
									jsii.String("/app/gossip"),
									jsii.String("--server"),
								},
							},
						},
					},
				},
			},
		},
	)

	k8s.NewKubeDeployment(
		chart.Cdk8sChart,
		jsii.String("deploy-client"),
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
								Env:             &env,
								Name:            jsii.String(appName),
								Image:           jsii.String(imageName),
								ImagePullPolicy: jsii.String("Always"),
								Command: &[]*string{
									jsii.String("/app/gossip"),
									jsii.String("--client"),
								},
							},
						},
					},
				},
			},
		},
	)

	return chart
}

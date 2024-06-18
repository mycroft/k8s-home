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
			Name:  jsii.String("DELAY"),
			Value: jsii.String("30"),
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
									jsii.String("server"),
									jsii.String("--nats"),
									jsii.String("nats.nats:4222"),
									jsii.String("--temporal"),
									jsii.String("temporal-frontend-headless.temporal:7233"),
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
									jsii.String("client"),
									jsii.String("--nats"),
									jsii.String("nats.nats:4222"),
								},
							},
						},
					},
				},
			},
		},
	)

	replicas := float64(3.)

	k8s.NewKubeDeployment(
		chart.Cdk8sChart,
		jsii.String("deploy-worker"),
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
									jsii.String("worker"),
									jsii.String("--temporal"),
									jsii.String("temporal-frontend-headless.temporal:7233"),
								},
							},
						},
					},
				},
				Replicas: &replicas,
			},
		},
	)

	return chart
}

package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewEmojivotoChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "emojivoto"
	ingressName := "emojivoto.services.mkz.me"
	imageEmojiSvc := kubehelpers.RegisterDockerImage("docker.l5d.io/buoyantio/emojivoto-emoji-svc")
	imageWeb := kubehelpers.RegisterDockerImage("docker.l5d.io/buoyantio/emojivoto-web")
	imageVotingSvc := kubehelpers.RegisterDockerImage("docker.l5d.io/buoyantio/emojivoto-voting-svc")

	chart := builder.NewChart(namespace)

	k8s.NewKubeNamespace(
		chart.Cdk8sChart,
		jsii.String(fmt.Sprintf("ns-%s", namespace)),
		&k8s.KubeNamespaceProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String(namespace),
				Annotations: &map[string]*string{
					"linkerd.io/inject": jsii.String("enabled"),
				},
			},
		},
	)

	sa := []string{
		"emoji",
		"voting",
		"web",
	}

	for _, saName := range sa {
		k8s.NewKubeServiceAccount(
			chart.Cdk8sChart,
			jsii.String(fmt.Sprintf("%s-sa", saName)),
			&k8s.KubeServiceAccountProps{
				Metadata: &k8s.ObjectMeta{
					Name:      jsii.String(saName),
					Namespace: jsii.String(namespace),
				},
			},
		)
	}

	deploymentLabels := map[string]*string{
		"app":     jsii.String("emoji"),
		"version": jsii.String("v11"),
	}
	k8s.NewKubeDeployment(
		chart.Cdk8sChart,
		jsii.String("emoji-deploy"),
		&k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Annotations: &map[string]*string{
					"linkerd.io/inject": jsii.String("enabled"),
				},
				Labels: &map[string]*string{
					"app.kubernetes.io/name":    jsii.String("emoji"),
					"app.kubernetes.io/part-of": jsii.String("emojivoto"),
					"app.kubernetes.io/version": jsii.String("v11"),
				},
				Name:      jsii.String("emoji"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.DeploymentSpec{
				Replicas: jsii.Number(1),
				Selector: &k8s.LabelSelector{
					MatchLabels: &deploymentLabels,
				},
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels: &deploymentLabels,
					},
					Spec: &k8s.PodSpec{
						ServiceAccountName: jsii.String("emoji"),
						Containers: &[]*k8s.Container{
							{
								Env: &[]*k8s.EnvVar{
									{Name: jsii.String("GRPC_PORT"), Value: jsii.String("8080")},
									{Name: jsii.String("PROM_PORT"), Value: jsii.String("8081")},
								},
								Image: jsii.String(imageEmojiSvc),
								Name:  jsii.String("emoji"),
								Ports: &[]*k8s.ContainerPort{
									{ContainerPort: jsii.Number(8080), Name: jsii.String("grpc")},
									{ContainerPort: jsii.Number(8081), Name: jsii.String("prom")},
								},
							},
						},
					},
				},
			},
		},
	)

	deploymentLabels = map[string]*string{
		"app":     jsii.String("vote-bot"),
		"version": jsii.String("v11"),
	}
	k8s.NewKubeDeployment(
		chart.Cdk8sChart,
		jsii.String("vote-bot-deploy"),
		&k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Annotations: &map[string]*string{
					"linkerd.io/inject": jsii.String("enabled"),
				},
				Labels: &map[string]*string{
					"app.kubernetes.io/name":    jsii.String("vote-bot"),
					"app.kubernetes.io/part-of": jsii.String("emojivoto"),
					"app.kubernetes.io/version": jsii.String("v11"),
				},
				Name:      jsii.String("vote-bot"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.DeploymentSpec{
				Replicas: jsii.Number(1),
				Selector: &k8s.LabelSelector{
					MatchLabels: &deploymentLabels,
				},
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels: &deploymentLabels,
					},
					Spec: &k8s.PodSpec{
						Containers: &[]*k8s.Container{
							{
								Command: &[]*string{
									jsii.String("emojivoto-vote-bot"),
								},
								Name: jsii.String("vote-bot"),
								Env: &[]*k8s.EnvVar{
									{Name: jsii.String("WEB_HOST"), Value: jsii.String("web.emojivoto:80")},
								},
								Image: jsii.String(imageWeb),
								Ports: &[]*k8s.ContainerPort{
									{ContainerPort: jsii.Number(8080), Name: jsii.String("grpc")},
									{ContainerPort: jsii.Number(8081), Name: jsii.String("prom")},
								},
							},
						},
					},
				},
			},
		},
	)

	deploymentLabels = map[string]*string{
		"app":     jsii.String("voting"),
		"version": jsii.String("v11"),
	}
	k8s.NewKubeDeployment(
		chart.Cdk8sChart,
		jsii.String("voting-deploy"),
		&k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Annotations: &map[string]*string{
					"linkerd.io/inject": jsii.String("enabled"),
				},
				Labels: &map[string]*string{
					"app.kubernetes.io/name":    jsii.String("voting"),
					"app.kubernetes.io/part-of": jsii.String("emojivoto"),
					"app.kubernetes.io/version": jsii.String("v11"),
				},
				Name:      jsii.String("voting"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.DeploymentSpec{
				Replicas: jsii.Number(1),
				Selector: &k8s.LabelSelector{
					MatchLabels: &deploymentLabels,
				},
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels: &deploymentLabels,
					},
					Spec: &k8s.PodSpec{
						Containers: &[]*k8s.Container{
							{
								Env: &[]*k8s.EnvVar{
									{Name: jsii.String("GRPC_PORT"), Value: jsii.String("8080")},
									{Name: jsii.String("PROM_PORT"), Value: jsii.String("8081")},
								},
								Image: jsii.String(imageVotingSvc),
								Name:  jsii.String("voting"),
								Ports: &[]*k8s.ContainerPort{
									{ContainerPort: jsii.Number(8080), Name: jsii.String("grpc")},
									{ContainerPort: jsii.Number(8081), Name: jsii.String("prom")},
								},
							},
						},
						ServiceAccountName: jsii.String("voting"),
					},
				},
			},
		},
	)

	deploymentLabels = map[string]*string{
		"app":     jsii.String("web"),
		"version": jsii.String("v11"),
	}
	k8s.NewKubeDeployment(
		chart.Cdk8sChart,
		jsii.String("web-deploy"),
		&k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Annotations: &map[string]*string{
					"linkerd.io/inject": jsii.String("enabled"),
				},
				Labels: &map[string]*string{
					"app.kubernetes.io/name":    jsii.String("web"),
					"app.kubernetes.io/part-of": jsii.String("emojivoto"),
					"app.kubernetes.io/version": jsii.String("v11"),
				},
				Name:      jsii.String("web"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.DeploymentSpec{
				Replicas: jsii.Number(1),
				Selector: &k8s.LabelSelector{
					MatchLabels: &deploymentLabels,
				},
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels: &deploymentLabels,
					},
					Spec: &k8s.PodSpec{
						Containers: &[]*k8s.Container{
							{
								Env: &[]*k8s.EnvVar{
									{Name: jsii.String("WEB_PORT"), Value: jsii.String("8080")},
									{Name: jsii.String("EMOJISVC_HOST"), Value: jsii.String("emoji.emojivoto:8080")},
									{Name: jsii.String("VOTINGSVC_HOST"), Value: jsii.String("voting.emojivoto:8080")},
									{Name: jsii.String("INDEX_BUNDLE"), Value: jsii.String("dist/index_bundle.js")},
								},
								Image: jsii.String(imageWeb),
								Name:  jsii.String("web"),
								Ports: &[]*k8s.ContainerPort{
									{ContainerPort: jsii.Number(8080), Name: jsii.String("http")},
								},
							},
						},
						ServiceAccountName: jsii.String("web"),
					},
				},
			},
		},
	)

	k8s.NewKubeService(
		chart.Cdk8sChart,
		jsii.String("emoji"),
		&k8s.KubeServiceProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("emoji"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.ServiceSpec{
				Ports: &[]*k8s.ServicePort{
					{Name: jsii.String("grpc"), Port: jsii.Number(8080), TargetPort: k8s.IntOrString_FromNumber(jsii.Number(8080))},
					{Name: jsii.String("prom"), Port: jsii.Number(8081), TargetPort: k8s.IntOrString_FromNumber(jsii.Number(8081))},
				},
				Selector: &map[string]*string{"app": jsii.String("emoji")},
			},
		},
	)

	k8s.NewKubeService(
		chart.Cdk8sChart,
		jsii.String("voting"),
		&k8s.KubeServiceProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("voting"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.ServiceSpec{
				Ports: &[]*k8s.ServicePort{
					{Name: jsii.String("grpc"), Port: jsii.Number(8080), TargetPort: k8s.IntOrString_FromNumber(jsii.Number(8080))},
					{Name: jsii.String("prom"), Port: jsii.Number(8081), TargetPort: k8s.IntOrString_FromNumber(jsii.Number(8081))},
				},
				Selector: &map[string]*string{"app": jsii.String("voting")},
			},
		},
	)

	svc := k8s.NewKubeService(
		chart.Cdk8sChart,
		jsii.String("web"),
		&k8s.KubeServiceProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("web"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.ServiceSpec{
				Ports: &[]*k8s.ServicePort{
					{Name: jsii.String("http"), Port: jsii.Number(80), TargetPort: k8s.IntOrString_FromNumber(jsii.Number(8080))},
				},
				Selector: &map[string]*string{"app": jsii.String("web")},
			},
		},
	)

	kubehelpers.NewAppIngress(
		builder.Context,
		chart.Cdk8sChart,
		map[string]*string{
			"app": jsii.String("web"),
		},
		namespace,
		8080,
		ingressName,
		*svc.Name(),
		map[string]string{},
	)

	return chart
}

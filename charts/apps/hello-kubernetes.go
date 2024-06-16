package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewHelloKubernetesChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "hello-kubernetes"
	image := kubehelpers.RegisterDockerImage("paulbouwer/hello-kubernetes")

	chart := builder.NewChart(appName)
	chart.NewNamespace(appName)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	env := []*k8s.EnvVar{
		{
			Name: jsii.String("KUBERNETES_NAMESPACE"),
			ValueFrom: &k8s.EnvVarSource{
				FieldRef: &k8s.ObjectFieldSelector{
					FieldPath: jsii.String("metadata.namespace"),
				},
			},
		},
		{
			Name: jsii.String("KUBERNETES_POD_NAME"),
			ValueFrom: &k8s.EnvVarSource{
				FieldRef: &k8s.ObjectFieldSelector{
					FieldPath: jsii.String("metadata.name"),
				},
			},
		},
		{
			Name: jsii.String("KUBERNETES_NODE_NAME"),
			ValueFrom: &k8s.EnvVarSource{
				FieldRef: &k8s.ObjectFieldSelector{
					FieldPath: jsii.String("spec.nodeName"),
				},
			},
		},
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
								Name:  jsii.String(appName),
								Image: jsii.String(image),
								Env:   &env,
							},
						},
					},
				},
			},
		},
	)

	// It is required to add the following annotation:
	// traefik.ingress.kubernetes.io/router.middlewares: traefik-forward-auth-traefik-forward-auth@kubernetescrd

	kubehelpers.NewAppIngress(
		builder.Context,
		chart.Cdk8sChart,
		labels,
		appName,
		8080,
		"hello-kubernetes.services.mkz.me",
		"",
		map[string]string{
			"traefik.ingress.kubernetes.io/router.middlewares": "traefik-forward-auth-traefik-forward-auth@kubernetescrd",
		},
	)

	return chart
}

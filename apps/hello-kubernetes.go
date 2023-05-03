package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewHelloKubernetesChart(scope constructs.Construct) cdk8s.Chart {
	appName := "hello-kubernetes"
	image := k8s_helpers.RegisterDockerImage("paulbouwer/hello-kubernetes:1.10.1")

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	k8s.NewKubeNamespace(
		chart,
		jsii.String("ns"),
		&k8s.KubeNamespaceProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String(appName),
			},
		},
	)

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

	k8s_helpers.NewAppIngress(
		chart,
		labels,
		appName,
		8080,
		"hello-kubernetes.services.mkz.me",
		map[string]string{
			"traefik.ingress.kubernetes.io/router.middlewares": "traefik-forward-auth-traefik-forward-auth@kubernetescrd",
		},
	)

	return chart
}

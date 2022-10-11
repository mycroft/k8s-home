package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewHelloKubernetesChart(scope constructs.Construct) cdk8s.Chart {
	appName := "hello-kubernetes"

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
								Image: jsii.String("paulbouwer/hello-kubernetes:1.10.1"),
							},
						},
					},
				},
			},
		},
	)

	svc := k8s.NewKubeService(
		chart,
		jsii.String("svc"),
		&k8s.KubeServiceProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(appName),
			},
			Spec: &k8s.ServiceSpec{
				Ports: &[]*k8s.ServicePort{
					{
						Name: jsii.String("http"),
						Port: jsii.Number(8080),
					},
				},
				Selector: &labels,
			},
		},
	)

	annotations := map[string]*string{
		"kubernetes.io/ingress.class":                        jsii.String("traefik"),
		"cert-manager.io/cluster-issuer":                     jsii.String("letsencrypt-staging"),
		"traefik.ingress.kubernetes.io/redirect-entry-point": jsii.String("https"),
		"traefik.ingress.kubernetes.io/redirect-permanent":   jsii.String("true"),
	}

	ingressHost := "hello-kubernetes.services.mkz.me"

	k8s.NewKubeIngress(
		chart,
		jsii.String("ingress"),
		&k8s.KubeIngressProps{
			Metadata: &k8s.ObjectMeta{
				Annotations: &annotations,
				Namespace:   jsii.String(appName),
			},
			Spec: &k8s.IngressSpec{
				Rules: &[]*k8s.IngressRule{
					{
						Host: jsii.String(ingressHost),
						Http: &k8s.HttpIngressRuleValue{
							Paths: &[]*k8s.HttpIngressPath{
								{
									Backend: &k8s.IngressBackend{
										Service: &k8s.IngressServiceBackend{
											Name: svc.Name(),
											Port: &k8s.ServiceBackendPort{
												Name: jsii.String("http"),
											},
										},
									},
									Path:     jsii.String("/"),
									PathType: jsii.String("Prefix"),
								},
							},
						},
					},
				},
				Tls: &[]*k8s.IngressTls{
					{
						Hosts: &[]*string{
							jsii.String(ingressHost),
						},
						SecretName: jsii.String("secret-tls-www"),
					},
				},
			},
		},
	)

	return chart
}

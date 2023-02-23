package k8s_helpers

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewAppIngress(
	chart cdk8s.Chart,
	labelSelector map[string]*string,
	appName string,
	appPort int,
	ingressHost string,
) {
	annotations := map[string]*string{
		"kubernetes.io/ingress.class":                        jsii.String("traefik"),
		"cert-manager.io/cluster-issuer":                     jsii.String("letsencrypt-prod"),
		"traefik.ingress.kubernetes.io/redirect-entry-point": jsii.String("https"),
		"traefik.ingress.kubernetes.io/redirect-permanent":   jsii.String("true"),
	}

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
						Port: jsii.Number(float64(appPort)),
					},
				},
				Selector: &labelSelector,
			},
		},
	)

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

}
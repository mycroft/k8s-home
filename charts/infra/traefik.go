package infra

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewTraefikChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "kube-system"
	ingressHost := "traefik.services.mkz.me"
	portName := "web"
	appPort := 9000

	annotations := map[string]*string{
		"cert-manager.io/cluster-issuer":                     jsii.String("letsencrypt-prod"),
		"traefik.ingress.kubernetes.io/router.middlewares":   jsii.String("traefik-forward-auth-traefik-forward-auth@kubernetescrd"),
		"traefik.ingress.kubernetes.io/redirect-entry-point": jsii.String("https"),
		"traefik.ingress.kubernetes.io/redirect-permanent":   jsii.String("true"),
		"traefik.ingress.kubernetes.io/app-root":             jsii.String("/dashboard/"),
	}

	chart := builder.NewChart("traefik")

	labels := map[string]*string{
		"app.kubernetes.io/instance": jsii.String("traefik-kube-system"),
		"app.kubernetes.io/name":     jsii.String("traefik"),
	}

	svc := kubehelpers.NewAppService(
		chart.Cdk8sChart,
		namespace,
		"svc",
		labels,
		portName,
		uint(appPort),
	)

	rules := []*k8s.IngressRule{}
	hosts := []*string{&ingressHost}

	rules = append(rules, &k8s.IngressRule{
		Host: jsii.String(ingressHost),
		Http: &k8s.HttpIngressRuleValue{
			Paths: &[]*k8s.HttpIngressPath{
				{
					Backend: &k8s.IngressBackend{
						Service: &k8s.IngressServiceBackend{
							Name: svc.Name(),
							Port: &k8s.ServiceBackendPort{
								Name: jsii.String(portName),
							},
						},
					},
					Path:     jsii.String("/"),
					PathType: jsii.String("Prefix"),
				},
			},
		},
	})

	k8s.NewKubeIngress(
		chart.Cdk8sChart,
		jsii.String("ingress"),
		&k8s.KubeIngressProps{
			Metadata: &k8s.ObjectMeta{
				Annotations: &annotations,
				Namespace:   jsii.String(namespace),
			},
			Spec: &k8s.IngressSpec{
				IngressClassName: jsii.String("traefik"),
				Rules:            &rules,
				Tls: &[]*k8s.IngressTls{
					{
						Hosts:      &hosts,
						SecretName: jsii.String("secret-tls-www"),
					},
				},
			},
		},
	)

	return chart
}

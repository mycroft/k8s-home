package security

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/imports/traefikio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewTraefikForwardAuth(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "traefik-forward-auth"
	namespace := appName
	image := kubehelpers.RegisterDockerImage("thomseddon/traefik-forward-auth")

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	domainName := "mkz.me"
	ingressHost := fmt.Sprintf("forward-auth.services.%s", domainName)

	env := &[]*k8s.EnvVar{
		{
			Name:  jsii.String("DEFAULT_PROVIDER"),
			Value: jsii.String("oidc"),
		},
		{
			Name:  jsii.String("AUTH_HOST"),
			Value: jsii.String(ingressHost),
		},
		{
			Name:  jsii.String("COOKIE_DOMAIN"),
			Value: jsii.String(domainName),
		},
		{
			Name: jsii.String("SECRET"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("secret"),
					Key:  jsii.String("secret"),
				},
			},
		},
		{
			Name: jsii.String("PROVIDERS_OIDC_CLIENT_ID"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("oidc"),
					Key:  jsii.String("client_id"),
				},
			},
		},
		{
			Name: jsii.String("PROVIDERS_OIDC_CLIENT_SECRET"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("oidc"),
					Key:  jsii.String("client_secret"),
				},
			},
		},
		{
			Name: jsii.String("PROVIDERS_OIDC_ISSUER_URL"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("oidc"),
					Key:  jsii.String("issuer_url"),
				},
			},
		},
	}

	// External secret for client_id/client_secret/issuer_url
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "oidc")
	// External secret for secret
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "secret")

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
								Env:   env,
							},
						},
					},
				},
			},
		},
	)

	svc := k8s.NewKubeService(
		chart.Cdk8sChart,
		jsii.String("svc"),
		&k8s.KubeServiceProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(appName),
			},
			Spec: &k8s.ServiceSpec{
				Ports: &[]*k8s.ServicePort{
					{
						Name: jsii.String("http"),
						Port: jsii.Number(4181),
					},
				},
				Selector: &labels,
			},
		},
	)

	traefikio.NewMiddleware(
		chart.Cdk8sChart,
		jsii.String("traefik-forward-auth-middleware"),
		&traefikio.MiddlewareProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("traefik-forward-auth"),
				Namespace: jsii.String(namespace),
			},
			Spec: &traefikio.MiddlewareSpec{
				ForwardAuth: &traefikio.MiddlewareSpecForwardAuth{
					Address:            jsii.String(fmt.Sprintf("http://%s.%s:4181", *svc.Name(), namespace)),
					TrustForwardHeader: jsii.Bool(true),
					AuthRequestHeaders: &[]*string{
						jsii.String("X-Forwarded-User"),
						jsii.String("Cookie"),
					},
				},
			},
		},
	)

	annotations := map[string]*string{
		"cert-manager.io/cluster-issuer":                     jsii.String("letsencrypt-prod"),
		"traefik.ingress.kubernetes.io/redirect-entry-point": jsii.String("https"),
		"traefik.ingress.kubernetes.io/redirect-permanent":   jsii.String("true"),
		// Do not remove this! It breaks everything.
		"traefik.ingress.kubernetes.io/router.middlewares": jsii.String("traefik-forward-auth-traefik-forward-auth@kubernetescrd"),
	}

	k8s.NewKubeIngress(
		chart.Cdk8sChart,
		jsii.String("ingress"),
		&k8s.KubeIngressProps{
			Metadata: &k8s.ObjectMeta{
				Annotations: &annotations,
				Namespace:   jsii.String(appName),
			},
			Spec: &k8s.IngressSpec{
				IngressClassName: jsii.String("traefik"),
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

//
// Set the following label to service to protect:
// traefik.ingress.kubernetes.io/router.middlewares: traefik-forward-auth-traefik-forward-auth@kubernetescrd
//

package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewTraefikForwardAuth(scope constructs.Construct) cdk8s.Chart {
	appName := "traefik-forward-auth"
	namespace := appName

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
				Name: jsii.String(namespace),
			},
		},
	)

	k8s_helpers.CreateSecretStore(chart, namespace)

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
	k8s_helpers.CreateExternalSecret(chart, namespace, "oidc")
	// External secret for secret
	k8s_helpers.CreateExternalSecret(chart, namespace, "secret")

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
								Image: jsii.String("thomseddon/traefik-forward-auth:2.2.0"),
								Env:   env,
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
						Port: jsii.Number(4181),
					},
				},
				Selector: &labels,
			},
		},
	)

	annotations := map[string]*string{
		"kubernetes.io/ingress.class":                        jsii.String("traefik"),
		"cert-manager.io/cluster-issuer":                     jsii.String("letsencrypt-prod"),
		"traefik.ingress.kubernetes.io/redirect-entry-point": jsii.String("https"),
		"traefik.ingress.kubernetes.io/redirect-permanent":   jsii.String("true"),
	}

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

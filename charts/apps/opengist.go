package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/imports/traefikio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewOpengistChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "opengist"

	namespace := appName
	appImage := builder.RegisterContainerImage("ghcr.io/thomiceli/opengist")
	appPort := 6157
	appIngress := "opengist.services.mkz.me"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "sso")

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	env := []*k8s.EnvVar{
		{
			Name: jsii.String("OG_OIDC_CLIENT_KEY"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("sso"),
					Key:  jsii.String("client_id"),
				},
			},
		},
		{
			Name: jsii.String("OG_OIDC_SECRET"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("sso"),
					Key:  jsii.String("client_secret"),
				},
			},
		},
		{
			Name: jsii.String("OG_OIDC_DISCOVERY_URL"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("sso"),
					Key:  jsii.String("discovery_url"),
				},
			},
		},
		{
			Name:  jsii.String("OG_SSH_PORT"),
			Value: jsii.String("22222"),
		},
	}

	_, serviceName := kubehelpers.NewStatefulSet(
		chart.Cdk8sChart,
		namespace,
		appName,
		appImage,
		appPort,
		labels,
		env,
		[]string{},
		[]kubehelpers.StatefulSetVolume{
			{
				Name:        "opengist",
				MountPath:   "/opengist",
				StorageSize: "8Gi",
			},
		},
	)

	kubehelpers.NewAppIngress(
		builder.Context,
		chart.Cdk8sChart,
		labels,
		appName,
		appPort,
		appIngress,
		serviceName,
		map[string]string{},
	)

	kubehelpers.NewAppService(
		chart.Cdk8sChart,
		namespace,
		"opengist-ssh-svc",
		labels,
		"opengist-ssh",
		uint(22222),
	)

	traefikio.NewIngressRouteTcp(
		chart.Cdk8sChart,
		jsii.String("ingress-route-tcp-opengist"),
		&traefikio.IngressRouteTcpProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("opengist-ssh"),
				Namespace: jsii.String(namespace),
			},
			Spec: &traefikio.IngressRouteTcpSpec{
				EntryPoints: &[]*string{
					jsii.String("opengist-ssh"),
				},
				Routes: &[]*traefikio.IngressRouteTcpSpecRoutes{
					{
						Match: jsii.String("HostSNI(`*`)"),
						Services: &[]*traefikio.IngressRouteTcpSpecRoutesServices{
							{
								// Namespace: jsii.String(namespace),
								Name: jsii.String("opengist-opengist-ssh-svc-c8db8f7c"),
								Port: traefikio.IngressRouteTcpSpecRoutesServicesPort_FromNumber(jsii.Number(22222)),
							},
						},
					},
				},
			},
		},
	)

	return chart
}

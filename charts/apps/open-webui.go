package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewOpenWebuiChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "open-webui"
	namespace := appName

	appIngress := "ai.services.mkz.me"
	appImage := builder.RegisterContainerImage("ghcr.io/open-webui/open-webui")
	appPort := uint(8080)

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "sso")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "openai")

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	env := []*k8s.EnvVar{
		{
			Name: jsii.String("OAUTH_CLIENT_ID"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("sso"),
					Key:  jsii.String("client_id"),
				},
			},
		},
		{
			Name: jsii.String("OAUTH_CLIENT_SECRET"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("sso"),
					Key:  jsii.String("client_secret"),
				},
			},
		},
		{
			Name: jsii.String("OPENID_PROVIDER_URL"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("sso"),
					Key:  jsii.String("discovery_url"),
				},
			},
		},
		{
			Name:  jsii.String("ENABLE_OAUTH_SIGNUP"),
			Value: jsii.String("true"),
		},
		{
			Name:  jsii.String("OLLAMA_BASE_URL"),
			Value: jsii.String("http://10.0.0.8:11434"),
		},
		{
			Name:  jsii.String("ENABLE_OPENAI_API"),
			Value: jsii.String("True"),
		},
		{
			Name: jsii.String("OPENAI_API_KEY"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("openai"),
					Key:  jsii.String("api_key"),
				},
			},
		},
		{
			Name:  jsii.String("ENABLE_LOGIN_FORM"),
			Value: jsii.String("False"),
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
		[]kubehelpers.ConfigMapMount{},
		[]kubehelpers.StatefulSetVolume{
			{
				Name:        "openai",
				MountPath:   "/openai",
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

	return chart
}

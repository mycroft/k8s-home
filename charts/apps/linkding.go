package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewLinkdingChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "linkding"
	appName := namespace
	appPort := 9090
	appIngress := "links.services.mkz.me"
	linkdingImage := builder.RegisterContainerImage("sissbruecker/linkding")

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "openid")

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	env := []*k8s.EnvVar{
		{
			Name: jsii.String("OIDC_RP_CLIENT_ID"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("openid"),
					Key:  jsii.String("client_id"),
				},
			},
		},
		{
			Name: jsii.String("OIDC_RP_CLIENT_SECRET"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("openid"),
					Key:  jsii.String("client_secret"),
				},
			},
		},
		{
			Name:  jsii.String("OIDC_OP_AUTHORIZATION_ENDPOINT"),
			Value: jsii.Sprintf("https://auth.services.mkz.me/application/o/authorize/"),
		},
		{
			Name:  jsii.String("OIDC_OP_TOKEN_ENDPOINT"),
			Value: jsii.Sprintf("https://auth.services.mkz.me/application/o/token/"),
		},
		{
			Name:  jsii.String("OIDC_OP_USER_ENDPOINT"),
			Value: jsii.Sprintf("https://auth.services.mkz.me/application/o/userinfo/"),
		},
		{
			Name:  jsii.String("OIDC_OP_JWKS_ENDPOINT"),
			Value: jsii.Sprintf("https://auth.services.mkz.me/application/o/linkding/jwks/"),
		},
		{
			Name:  jsii.String("LD_ENABLE_OIDC"),
			Value: jsii.Sprintf("True"),
		},
	}

	_, svcName := kubehelpers.NewStatefulSet(
		chart.Cdk8sChart,
		namespace,
		appName,
		linkdingImage,
		appPort,
		labels,
		env,
		[]string{},
		[]kubehelpers.ConfigMapMount{},
		[]kubehelpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/etc/linkding/data",
				StorageSize: "1Gi",
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
		svcName,
		map[string]string{},
	)

	return chart
}

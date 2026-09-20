package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewVaultWardenChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "vaultwarden"
	appName := namespace
	appImage := builder.RegisterContainerImage("vaultwarden/server")
	appPort := uint(80)
	appIngress := "vaultwarden.services.mkz.me"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)
	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "admin")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "mail")

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	env := []*k8s.EnvVar{
		{
			Name:  jsii.String("DOMAIN"),
			Value: jsii.String("https://" + appIngress),
		},
		{
			Name: jsii.String("ADMIN_TOKEN"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("admin"),
					Key:  jsii.String("token"),
				},
			},
		},
		{
			Name: jsii.String("SMTP_HOST"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("mail"),
					Key:  jsii.String("hostname"),
				},
			},
		},
		{
			Name: jsii.String("SMTP_USERNAME"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("mail"),
					Key:  jsii.String("username"),
				},
			},
		},
		{
			Name: jsii.String("SMTP_PASSWORD"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("mail"),
					Key:  jsii.String("password"),
				},
			},
		},
		{Name: jsii.String("SMTP_PORT"), Value: jsii.String("587")},
		{Name: jsii.String("SMTP_SECURITY"), Value: jsii.String("starttls")},
		{Name: jsii.String("SMTP_FROM"), Value: jsii.String("vaultwarden@mkz.me")},
	}

	_, svcName := kubehelpers.NewStatefulSet(chart.Cdk8sChart, kubehelpers.StatefulSetConfig{
		Namespace: namespace,
		AppName:   appName,
		AppImage:  appImage,
		AppPort:   appPort,
		Labels:    labels,
		Env:       env,
		Storages: []kubehelpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/data",
				StorageSize: "10Gi",
			},
		},
	})

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

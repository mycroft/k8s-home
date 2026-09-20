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

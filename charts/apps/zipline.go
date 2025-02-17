package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewZiplineChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "zipline"
	appName := namespace
	appImage := builder.RegisterContainerImage("ghcr.io/diced/zipline")
	appPort := 3000
	appIngress := "zipline.services.mkz.me"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql_v4")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "core")

	env := []*k8s.EnvVar{
		{
			Name: jsii.String("CORE_SECRET"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Key:  jsii.String("secret"),
					Name: jsii.String("core"),
				},
			},
		},
		{
			Name: jsii.String("DATABASE_URL"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Key:  jsii.String("url"),
					Name: jsii.String("postgresql_v4"),
				},
			},
		},
		{
			Name: jsii.String("CORE_DATABASE_URL"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Key:  jsii.String("url"),
					Name: jsii.String("postgresql"),
				},
			},
		},
	}

	_, svcName := kubehelpers.NewStatefulSet(
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
				Name:        "zipline-uploads",
				MountPath:   "/zipline/uploads",
				StorageSize: "32Gi",
			},
			{
				Name:        "zipline-public",
				MountPath:   "/zipline/public",
				StorageSize: "32Gi",
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

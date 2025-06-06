package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewVikunjaChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "vikunja"

	namespace := appName
	appImage := builder.RegisterContainerImage("vikunja/vikunja")
	appPort := uint(3456)
	appIngress := "vikunja.services.mkz.me"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql")
	// kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "openid")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "config") // key: config.yml

	env := []*k8s.EnvVar{
		{
			Name:  jsii.String("VIKUNJA_DATABASE_TYPE"),
			Value: jsii.String("postgres"),
		},
		{
			Name:  jsii.String("VIKUNJA_DATABASE_NAME"),
			Value: jsii.String("vikunja"),
		},
		{
			Name:  jsii.String("VIKUNJA_DATABASE_HOST"),
			Value: jsii.String("postgres-instance.postgres"),
		},
		{
			Name: jsii.String("VIKUNJA_DATABASE_USER"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Key:  jsii.String("username"),
					Name: jsii.String("postgresql"),
				},
			},
		},
		{
			Name: jsii.String("VIKUNJA_DATABASE_PASSWORD"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Key:  jsii.String("password"),
					Name: jsii.String("postgresql"),
				},
			},
		},
	}

	// config, err := os.ReadFile("configs/vikunja/config.yaml")
	// if err != nil {
	// 	log.Fatalf("Could not read privatebin configuration file: %s", err.Error())
	// }

	// configMap := k8s.NewKubeConfigMap(
	// 	chart.Cdk8sChart,
	// 	jsii.String("config"),
	// 	&k8s.KubeConfigMapProps{
	// 		Metadata: &k8s.ObjectMeta{
	// 			Namespace: jsii.String(namespace),
	// 		},
	// 		Data: &map[string]*string{
	// 			"config.yml": jsii.String(string(config)),
	// 		},
	// 	},
	// )

	secrets := []kubehelpers.SecretMount{
		{
			Name:      "config",
			MountPath: "/etc/vikunja",
		},
	}

	_, svcName := kubehelpers.NewStatefulSetWithSecrets(
		chart.Cdk8sChart,
		namespace,
		appName,
		appImage,
		appPort,
		labels,
		env,
		[]string{},
		[]kubehelpers.ConfigMapMount{
			// {
			// 	Name:      "config",
			// 	ConfigMap: configMap,
			// 	MountPath: "/etc/vikunja",
			// },
		},
		secrets,
		[]kubehelpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/app/vikunja/files",
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

package apps

import (
	"strings"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewVikunjaChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "vikunja"

	namespace := appName
	appImage := builder.RegisterContainerImage("vikunja/vikunja")
	appPort := uint(3456)
	appIngresses := []string{
		"vikunja.services.mkz.me",
		"todo.iop.cx",
	}

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql-cnpg")
	// kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "openid")
	// Vikunja's config.yml is managed in Vault, mounted from the `config` secret below.
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "config")

	// Vikunja appends service.publicurl to cors.origins on its own, but every
	// other hostname has to be listed or the browser's XHRs get no
	// Access-Control-Allow-Origin. Setting this replaces the upstream default,
	// so the localhost entries the desktop app relies on are repeated here.
	// viper reads list values from the environment space-separated.
	corsOrigins := []string{"http://127.0.0.1:*", "http://localhost:*"}
	for _, host := range appIngresses {
		corsOrigins = append(corsOrigins, "https://"+host)
	}

	env := []*k8s.EnvVar{
		{
			Name:  jsii.String("VIKUNJA_CORS_ORIGINS"),
			Value: jsii.String(strings.Join(corsOrigins, " ")),
		},
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
			Value: jsii.String("postgres-rw.cnpg"),
		},
		{
			Name: jsii.String("VIKUNJA_DATABASE_USER"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Key:  jsii.String("username"),
					Name: jsii.String("postgresql-cnpg"),
				},
			},
		},
		{
			Name: jsii.String("VIKUNJA_DATABASE_PASSWORD"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Key:  jsii.String("password"),
					Name: jsii.String("postgresql-cnpg"),
				},
			},
		},
	}

	secrets := []kubehelpers.SecretMount{
		{
			Name:      "config",
			MountPath: "/etc/vikunja",
		},
	}

	_, svcName := kubehelpers.NewStatefulSet(chart.Cdk8sChart, kubehelpers.StatefulSetConfig{
		Namespace:    namespace,
		AppName:      appName,
		AppImage:     appImage,
		AppPort:      appPort,
		Labels:       labels,
		Env:          env,
		SecretMounts: secrets,
		Storages: []kubehelpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/app/vikunja/files",
				StorageSize: "1Gi",
			},
		},
		FsGroup: 1000,
	})

	kubehelpers.NewAppIngresses(
		builder.Context,
		chart.Cdk8sChart,
		labels,
		appName,
		appPort,
		appIngresses,
		svcName,
		map[string]string{},
	)

	return chart
}

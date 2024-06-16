package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewPaperlessNGXChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "paperless-ngx"
	appIngress := "paperless.services.mkz.me"

	paperlessNgxImage := builder.RegisterContainerImage("paperlessngx/paperless-ngx")

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql")

	_, redisServiceName := chart.NewRedisStatefulset(namespace)

	env := []*k8s.EnvVar{
		// XXX fix url here
		{
			Name:  jsii.String("PAPERLESS_REDIS"),
			Value: jsii.String(fmt.Sprintf("redis://%s:6379", redisServiceName)),
		},
		{Name: jsii.String("PAPERLESS_DBENGINE"), Value: jsii.String("postgresql")},
		{Name: jsii.String("PAPERLESS_DBHOST"), Value: jsii.String("postgres-instance.postgres")},
		{Name: jsii.String("PAPERLESS_DBPORT"), Value: jsii.String("5432")},
		{Name: jsii.String("PAPERLESS_DBNAME"), Value: jsii.String("paperlessngx")},
		{
			Name: jsii.String("PAPERLESS_DBUSER"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("postgresql"),
					Key:  jsii.String("username"),
				},
			},
		},
		{
			Name: jsii.String("PAPERLESS_DBPASS"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("postgresql"),
					Key:  jsii.String("password"),
				},
			},
		},
		{
			Name:  jsii.String("PAPERLESS_URL"),
			Value: jsii.String(fmt.Sprintf("https://%s", appIngress)),
		},
		{
			Name:  jsii.String("PAPERLESS_MEDIA_ROOT"),
			Value: jsii.String("/usr/src/paperless/media"),
		},
		{
			Name:  jsii.String("PAPERLESS_DATA_DIR"),
			Value: jsii.String("/usr/src/paperless/data"),
		},
	}

	paperlessngxLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("paperless-ngx"),
	}

	appName := "paperless-ngx"
	appPort := 8000

	_, svcName := kubehelpers.NewStatefulSet(
		chart.Cdk8sChart,
		namespace,
		appName,
		paperlessNgxImage,
		appPort,
		paperlessngxLabels,
		env,
		[]string{},
		[]kubehelpers.StatefulSetVolume{
			{ // PAPERLESS_DATA_DIR
				Name:        "data",
				MountPath:   "/usr/src/paperless/data",
				StorageSize: "2Gi",
			},
			{ // PAPERLESS_MEDIA_ROOT
				Name:        "media",
				MountPath:   "/usr/src/paperless/media",
				StorageSize: "32Gi",
			},
		},
	)

	kubehelpers.NewAppIngress(
		builder.Context,
		chart.Cdk8sChart,
		paperlessngxLabels,
		appName,
		appPort,
		appIngress,
		svcName,
		map[string]string{},
	)

	return chart
}

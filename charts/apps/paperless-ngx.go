package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewPaperlessNGXChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "paperless-ngx"
	appIngress := "paperless.services.mkz.me"

	redisImage := k8s_helpers.RegisterDockerImage("redis")
	paperlessNgxImage := k8s_helpers.RegisterDockerImage("paperlessngx/paperless-ngx")

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)
	k8s_helpers.CreateExternalSecret(chart, namespace, "postgresql")

	redisLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("redis"),
	}

	k8s_helpers.NewStatefulSet(
		chart,
		namespace,
		"redis",
		redisImage,
		6379,
		redisLabels,
		[]*k8s.EnvVar{},
		[]string{
			"redis-server --save 60 1 --loglevel warning",
		},
		[]k8s_helpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/data",
				StorageSize: "1Gi",
			},
		},
	)

	env := []*k8s.EnvVar{
		// XXX fix url here
		{
			Name:  jsii.String("PAPERLESS_REDIS"),
			Value: jsii.String("redis://paperless-ngx-redis-svc-c86d771a:6379"),
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

	k8s_helpers.NewStatefulSet(
		chart,
		namespace,
		appName,
		paperlessNgxImage,
		appPort,
		paperlessngxLabels,
		env,
		[]string{},
		[]k8s_helpers.StatefulSetVolume{
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

	k8s_helpers.NewAppIngress(
		chart,
		paperlessngxLabels,
		appName,
		appPort,
		appIngress,
		map[string]string{},
	)

	return chart
}

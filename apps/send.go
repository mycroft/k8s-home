package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewSendChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "send"
	appName := "send"
	ingressHost := "send.services.mkz.me"
	appPort := 1443
	image := k8s_helpers.RegisterDockerImage("registry.gitlab.com/timvisee/send:latest")

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)
	k8s_helpers.CreateExternalSecret(chart, namespace, "minio")

	redisLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("redis"),
	}
	sendLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("send"),
	}

	k8s_helpers.NewStatefulSet(
		chart,
		namespace,
		"redis",
		"redis:7.2.0",
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

	redisHost := "send-redis-svc-c8e30a3c.send"

	env := []*k8s.EnvVar{
		{Name: jsii.String("BASE_URL"), Value: jsii.String(ingressHost)},
		{Name: jsii.String("REDIS_HOST"), Value: jsii.String(redisHost)},
		{Name: jsii.String("S3_ENDPOINT"), Value: jsii.String("https://minio-storage.services.mkz.me")},
		{Name: jsii.String("S3_BUCKET"), Value: jsii.String("send")},
		{Name: jsii.String("S3_USE_PATH_STYLE_ENDPOINT"), Value: jsii.String("true")},
		{
			Name: jsii.String("AWS_ACCESS_KEY_ID"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("minio"),
					Key:  jsii.String("AWS_ACCESS_KEY_ID"),
				},
			},
		},
		{
			Name: jsii.String("AWS_SECRET_ACCESS_KEY"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("minio"),
					Key:  jsii.String("AWS_ACCESS_SECRET_KEY"),
				},
			},
		},
	}

	k8s_helpers.NewAppDeployment(
		chart,
		namespace,
		appName,
		image,
		sendLabels,
		env,
		[]string{},
		[]k8s_helpers.ConfigMapMount{},
	)

	k8s_helpers.NewAppIngress(
		chart,
		sendLabels,
		appName,
		appPort,
		ingressHost,
		map[string]string{},
	)

	return chart
}

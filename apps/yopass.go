package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewYopassChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "yopass"
	appIngress := "yopass.services.mkz.me"
	appName := "yopass"
	appPort := 1337

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	redisLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("redis"),
	}

	k8s_helpers.NewStatefulSet(
		chart,
		namespace,
		"redis",
		"redis:7.0.9",
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
		false,
	)

	redisUrl := "redis://yopass-redis-svc-c8a159bf:6379"

	yopassLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("paperless-ngx"),
	}

	k8s_helpers.NewAppDeployment(
		chart,
		namespace,
		appName,
		"jhaals/yopass:11.5.0",
		yopassLabels,
		[]*k8s.EnvVar{},
		[]string{
			fmt.Sprintf("/yopass-server --database redis --metrics-port 1338 --port 1337 --redis %s", redisUrl),
		},
		[]k8s_helpers.ConfigMapMount{},
	)

	k8s_helpers.NewAppIngress(
		chart,
		yopassLabels,
		appName,
		appPort,
		appIngress,
		map[string]string{},
	)

	return chart
}

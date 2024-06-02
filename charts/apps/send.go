package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewSendChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "send"
	appName := "send"
	ingressHost := "send.services.mkz.me"
	appPort := 1443
	image := kubehelpers.RegisterDockerImage("registry.gitlab.com/timvisee/send")

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)
	kubehelpers.CreateSecretStore(chart, namespace)
	kubehelpers.CreateExternalSecret(chart, namespace, "minio")

	sendLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("send"),
	}

	_, redisServiceName := kubehelpers.NewRedisStatefulset(chart, namespace)

	redisHost := fmt.Sprintf("%s.%s", redisServiceName, namespace)

	env := []*k8s.EnvVar{
		{Name: jsii.String("BASE_URL"), Value: jsii.String(fmt.Sprintf("https://%s", ingressHost))},
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

	kubehelpers.NewAppDeployment(
		chart,
		namespace,
		appName,
		image,
		sendLabels,
		env,
		[]string{},
		[]kubehelpers.ConfigMapMount{},
	)

	kubehelpers.NewAppIngress(
		chart,
		sendLabels,
		appName,
		appPort,
		ingressHost,
		"",
		map[string]string{},
	)

	return chart
}

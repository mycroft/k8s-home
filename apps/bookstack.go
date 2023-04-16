package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewBookstackChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "bookstack"
	appName := namespace
	appImage := "linuxserver/bookstack:23.02.3"
	appPort := 80
	appIngress := "bookstack.services.mkz.me"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)
	k8s_helpers.CreateExternalSecret(chart, namespace, "mariadb")

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	env := []*k8s.EnvVar{
		{Name: jsii.String("APP_URL"), Value: jsii.String(fmt.Sprintf("https://%s", appIngress))},
		{Name: jsii.String("PUID"), Value: jsii.String("1000")},
		{Name: jsii.String("PGID"), Value: jsii.String("1000")},
		{Name: jsii.String("TZ"), Value: jsii.String("Etc/UTC")},
		{Name: jsii.String("DB_HOST"), Value: jsii.String("mariadb.mariadb")},
		{Name: jsii.String("DB_PORT"), Value: jsii.String("3306")},
		{Name: jsii.String("DB_USER"), Value: jsii.String("bookstack")},
		{Name: jsii.String("DB_DATABASE"), Value: jsii.String("bookstack")},
		{
			Name: jsii.String("DB_PASS"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("mariadb"),
					Key:  jsii.String("password"),
				},
			},
		},
	}

	k8s_helpers.NewStatefulSet(
		chart,
		namespace,
		appName,
		appImage,
		appPort,
		labels,
		env,
		[]string{},
		[]k8s_helpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/config",
				StorageSize: "1Gi",
			},
		},
	)

	k8s_helpers.NewAppIngress(
		chart,
		labels,
		appName,
		appPort,
		appIngress,
		map[string]string{},
	)

	return chart
}

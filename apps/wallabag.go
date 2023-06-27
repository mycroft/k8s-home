package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewWallabagChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "wallabag"
	appName := namespace
	appImage := k8s_helpers.RegisterDockerImage("wallabag/wallabag:2.6.1")
	appIngress := "wallabag.services.mkz.me"
	appPort := 80

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	k8s_helpers.CreateExternalSecret(chart, namespace, "postgresql")
	k8s_helpers.CreateExternalSecret(chart, namespace, "mailrelay")

	env := []*k8s.EnvVar{
		{Name: jsii.String("SYMFONY__ENV__DATABASE_DRIVER"), Value: jsii.String("pdo_pgsql")},
		{Name: jsii.String("SYMFONY__ENV__DATABASE_HOST"), Value: jsii.String("postgres-instance.postgres")},
		{Name: jsii.String("SYMFONY__ENV__DATABASE_PORT"), Value: jsii.String("5432")},
		{Name: jsii.String("SYMFONY__ENV__DATABASE_NAME"), Value: jsii.String("wallabag")},
		{Name: jsii.String("POPULATE_DATABASE"), Value: jsii.String("True")},
		{
			Name: jsii.String("SYMFONY__ENV__DATABASE_USER"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("postgresql"),
					Key:  jsii.String("username"),
				},
			},
		},
		{
			Name: jsii.String("SYMFONY__ENV__DATABASE_PASSWORD"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("postgresql"),
					Key:  jsii.String("password"),
				},
			},
		},
		{Name: jsii.String("SYMFONY__ENV__DOMAIN_NAME"), Value: jsii.String(fmt.Sprintf("https://%s", appIngress))},
		{
			Name: jsii.String("SYMFONY__ENV__MAILER_USER"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("mailrelay"),
					Key:  jsii.String("username"),
				},
			},
		},
		{
			Name: jsii.String("SYMFONY__ENV__MAILER_PASSWORD"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("mailrelay"),
					Key:  jsii.String("password"),
				},
			},
		},
		{Name: jsii.String("SYMFONY__ENV__MAILER_TRANSPORT"), Value: jsii.String("smtp")},
		{Name: jsii.String("SYMFONY__ENV__MAILER_HOST"), Value: jsii.String("maki.mkz.me")},
		{Name: jsii.String("SYMFONY__ENV__MAILER_PORT"), Value: jsii.String("465")},
		{Name: jsii.String("SYMFONY__ENV__MAILER_ENCRYPTION"), Value: jsii.String("ssl")},
		{Name: jsii.String("SYMFONY__ENV__MAILER_AUTH_MODE"), Value: jsii.String("plain")},
		{Name: jsii.String("SYMFONY__ENV__FROM_EMAIL"), Value: jsii.String("wallabag@mkz.me")},
		{Name: jsii.String("SYMFONY__ENV__FOSUSER_REGISTRATION"), Value: jsii.String("false")},
	}

	k8s_helpers.NewAppDeployment(
		chart,
		namespace,
		appName,
		appImage,
		labels,
		env,
		[]string{},
		[]k8s_helpers.ConfigMapMount{},
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

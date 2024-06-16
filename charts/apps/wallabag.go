package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewWallabagChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "wallabag"
	appName := namespace
	appImage := builder.RegisterContainerImage("wallabag/wallabag")
	appIngress := "wallabag.services.mkz.me"
	appPort := 80

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "mailrelay")

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

	kubehelpers.NewAppDeployment(
		chart.Cdk8sChart,
		namespace,
		appName,
		appImage,
		labels,
		env,
		[]string{},
		[]kubehelpers.ConfigMapMount{},
	)

	kubehelpers.NewAppIngress(
		builder.Context,
		chart.Cdk8sChart,
		labels,
		appName,
		appPort,
		appIngress,
		"",
		map[string]string{},
	)

	return chart
}

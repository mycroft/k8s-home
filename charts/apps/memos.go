package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewMemosChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "memos"
	namespace := appName
	appPort := 5230
	appIngress := fmt.Sprintf("%s.services.mkz.me", appName)
	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}
	env := []*k8s.EnvVar{
		{
			Name:  jsii.String("MEMOS_DRIVER"),
			Value: jsii.String("postgres"),
		},
		{
			Name: jsii.String("MEMOS_DSN"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Key:  jsii.String("dsn"),
					Name: jsii.String("postgres"),
				},
			},
		},
	}

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)
	chart.CreateSecretStore(namespace)
	chart.CreateExternalSecret(namespace, "postgres")

	kubehelpers.NewAppDeployment(
		chart.Cdk8sChart,
		namespace,
		appName,
		chart.Builder.RegisterContainerImage("neosmemo/memos"),
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

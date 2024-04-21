package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewMicrobinChart(scope constructs.Construct) cdk8s.Chart {
	appName := "microbin"
	namespace := appName
	appIngress := "bin.iop.cx"

	// appImage := k8s_helpers.RegisterDockerImage("danielszabo99/microbin")
	appImage := "ghcr.io/zhaobenny/microbin:latest"
	appPort := 8080

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)
	k8s_helpers.CreateExternalSecret(chart, namespace, "admin")

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	env := []*k8s.EnvVar{
		{Name: jsii.String("MICROBIN_LIST_SERVER"), Value: jsii.Sprintf("FALSE")},
		{Name: jsii.String("MICROBIN_ENABLE_BURN_AFTER"), Value: jsii.Sprintf("TRUE")},
		{Name: jsii.String("MICROBIN_GC_DAYS"), Value: jsii.String("0")},
		{Name: jsii.String("MICROBIN_PRIVATE"), Value: jsii.String("TRUE")},
		{Name: jsii.Sprintf("MICROBIN_ENCRYPTION_CLIENT_SIDE"), Value: jsii.String("FALSE")},
		{Name: jsii.Sprintf("MICROBIN_ENCRYPTION_SERVER_SIDE"), Value: jsii.String("TRUE")},
		{Name: jsii.Sprintf("MICROBIN_WIDE"), Value: jsii.String("TRUE")},
		{Name: jsii.String("MICROBIN_PUBLIC_PATH"), Value: jsii.String(fmt.Sprintf("https://%s/", appIngress))},
		{
			Name: jsii.String("MICROBIN_ADMIN_USERNAME"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("admin"),
					Key:  jsii.String("username"),
				},
			},
		},
		{
			Name: jsii.String("MICROBIN_ADMIN_PASSWORD"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("admin"),
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
				Name:        "micro-bin",
				MountPath:   "/app/microbin_data",
				StorageSize: "4Gi",
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

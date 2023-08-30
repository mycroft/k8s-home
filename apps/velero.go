package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/veleroio"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewVeleroChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "velero"
	repositoryName := "vmware-tanzu"
	chartName := "velero"
	releaseName := "velero"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)
	k8s_helpers.CreateExternalSecret(chart, namespace, "nas0-minio")

	k8s_helpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://vmware-tanzu.github.io/helm-charts",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName,
		chartName,
		releaseName,
		"5.0.2",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"velero.yaml",
			),
		},
		nil,
	)

	// Create a default backup
	veleroio.NewSchedule(
		chart,
		jsii.String("backup-schedule"),
		&veleroio.ScheduleProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("backup-schedule"),
			},
			Spec: &veleroio.ScheduleSpec{
				Schedule: jsii.String("30 7 * * *"),
				Template: &veleroio.ScheduleSpecTemplate{
					Ttl: jsii.String("720h0m0s"),
				},
			},
		},
	)

	return chart
}

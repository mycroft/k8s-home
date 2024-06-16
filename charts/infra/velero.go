package infra

import (
	"git.mkz.me/mycroft/k8s-home/imports/veleroio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewVeleroChart(builder *kubehelpers.Builder) cdk8s.Chart {
	namespace := "velero"
	repositoryName := "vmware-tanzu"
	chartName := "velero"
	releaseName := "velero"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "nas0-minio")

	kubehelpers.CreateHelmRepository(
		chart.Cdk8sChart,
		repositoryName,
		"https://vmware-tanzu.github.io/helm-charts",
	)

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,
		repositoryName,
		chartName,
		releaseName,
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				releaseName,
				"velero.yaml",
			),
		},
		nil,
	)

	// Create a default backup
	veleroio.NewSchedule(
		chart.Cdk8sChart,
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

	return chart.Cdk8sChart
}

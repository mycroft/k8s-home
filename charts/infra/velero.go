package infra

import (
	"context"

	"git.mkz.me/mycroft/k8s-home/imports/veleroio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewVeleroChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	namespace := "velero"
	repositoryName := "vmware-tanzu"
	chartName := "velero"
	releaseName := "velero"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)
	kubehelpers.CreateSecretStore(chart, namespace)
	kubehelpers.CreateExternalSecret(chart, namespace, "nas0-minio")

	kubehelpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://vmware-tanzu.github.io/helm-charts",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName,
		chartName,
		releaseName,
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
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

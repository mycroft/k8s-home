package infra

import (
	"git.mkz.me/mycroft/k8s-home/imports/veleroio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewVeleroChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "velero"
	repositoryName := "vmware-tanzu"
	chartName := "velero"
	releaseName := "velero"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "garage-credentials")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "velero-repo-credentials")

	chart.CreateHelmRepository(
		repositoryName,
		"https://vmware-tanzu.github.io/helm-charts",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	postgresDumpHooks := []*veleroio.ScheduleSpecTemplateHooksResources{
		{
			Name:               jsii.String("postgres-dump-ready"),
			IncludedNamespaces: &[]*string{jsii.String("cnpg")},
			IncludedResources:  &[]*string{jsii.String("pods")},
			LabelSelector: &veleroio.ScheduleSpecTemplateHooksResourcesLabelSelector{
				MatchLabels: &map[string]*string{
					"app.kubernetes.io/name": jsii.String("postgres-dump"),
				},
			},
			Pre: &[]*veleroio.ScheduleSpecTemplateHooksResourcesPre{
				{
					Exec: &veleroio.ScheduleSpecTemplateHooksResourcesPreExec{
						Container: jsii.String("postgres-dump"),
						Command: &[]*string{
							jsii.String("/bin/sh"),
							jsii.String("-c"),
							jsii.String(`test -s /backup/latest/globals.sql && find -L /backup/latest -maxdepth 0 -type d -mmin -180 -print -quit | grep -q . && find -L /backup/latest -maxdepth 1 -type f -name '*.dump' -size +0c -print -quit | grep -q .`),
						},
						OnError: veleroio.ScheduleSpecTemplateHooksResourcesPreExecOnError_FAIL,
						Timeout: jsii.String("1m"),
					},
				},
			},
		},
	}

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

	// Back up every mounted persistent volume using Kopia. This includes the
	// timestamped logical PostgreSQL dumps produced in the cnpg namespace.
	veleroio.NewSchedule(
		chart.Cdk8sChart,
		jsii.String("filesystem-backup-schedule"),
		&veleroio.ScheduleProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("filesystem-backup-schedule"),
			},
			Spec: &veleroio.ScheduleSpec{
				Schedule: jsii.String("30 1 * * *"),
				Template: &veleroio.ScheduleSpecTemplate{
					DefaultVolumesToFsBackup: jsii.Bool(true),
					Hooks: &veleroio.ScheduleSpecTemplateHooks{
						Resources: &postgresDumpHooks,
					},
					SnapshotVolumes: jsii.Bool(false),
					Ttl:             jsii.String("720h0m0s"),
				},
			},
		},
	)

	return chart
}

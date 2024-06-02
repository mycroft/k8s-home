package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewFreshRSS(scope constructs.Construct) cdk8s.Chart {
	namespace := "freshrss"
	appName := namespace
	appImage := kubehelpers.RegisterDockerImage("freshrss/freshrss")
	appPort := 80
	appIngress := "freshrss.services.mkz.me"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	env := []*k8s.EnvVar{
		{Name: jsii.String("PUID"), Value: jsii.String("1000")},
		{Name: jsii.String("PGID"), Value: jsii.String("1000")},
		{Name: jsii.String("TZ"), Value: jsii.String("Etc/UTC")},
	}

	kubehelpers.NewNamespace(chart, namespace)

	stsName, _ := kubehelpers.NewStatefulSet(
		chart,
		namespace,
		appName,
		appImage,
		appPort,
		labels,
		env,
		[]string{},
		[]kubehelpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/var/www/FreshRSS/data",
				StorageSize: "1Gi",
			},
			{
				Name:        "extensions",
				MountPath:   "/var/www/FreshRSS/extensions",
				StorageSize: "1Gi",
			},
		},
	)

	affinity := &k8s.Affinity{
		PodAffinity: &k8s.PodAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &[]*k8s.PodAffinityTerm{
				{
					TopologyKey: jsii.String("kubernetes.io/hostname"),
					LabelSelector: &k8s.LabelSelector{
						MatchExpressions: &[]*k8s.LabelSelectorRequirement{
							{
								Key:      jsii.String("statefulset.kubernetes.io/pod-name"),
								Operator: jsii.String("In"),
								Values: &[]*string{
									jsii.String(fmt.Sprintf("%s-0", stsName)),
								},
							},
						},
					},
				},
			},
		},
	}

	k8s.NewKubeCronJob(
		chart,
		jsii.String("cronjob"),
		&k8s.KubeCronJobProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.CronJobSpec{
				Schedule: jsii.String("12,42 * * * *"),
				JobTemplate: &k8s.JobTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Namespace: jsii.String(namespace),
					},
					Spec: &k8s.JobSpec{
						Template: &k8s.PodTemplateSpec{
							Metadata: &k8s.ObjectMeta{
								Namespace: jsii.String(namespace),
							},
							Spec: &k8s.PodSpec{
								Affinity: affinity,
								Containers: &[]*k8s.Container{
									{
										Name: jsii.String("updater"),
										Command: &[]*string{
											jsii.String("/var/www/FreshRSS/app/actualize_script.php"),
										},
										Image: jsii.String(appImage),
										VolumeMounts: &[]*k8s.VolumeMount{
											{
												Name:      jsii.String("data"),
												MountPath: jsii.String("/var/www/FreshRSS/data"),
											},
										},
									},
								},
								RestartPolicy: jsii.String("Never"),
								Volumes: &[]*k8s.Volume{
									{
										Name: jsii.String("data"),
										PersistentVolumeClaim: &k8s.PersistentVolumeClaimVolumeSource{
											ClaimName: jsii.String(fmt.Sprintf("data-%s-0", stsName)),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	)

	kubehelpers.NewAppIngress(
		chart,
		labels,
		appName,
		appPort,
		appIngress,
		"",
		map[string]string{},
	)

	return chart
}

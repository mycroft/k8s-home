package apps

import (
	"slices"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewFreshRSS(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "freshrss"
	appName := namespace
	appImage := builder.RegisterContainerImage("freshrss/freshrss")
	appPort := uint(80)
	appIngress := "freshrss.services.mkz.me"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)
	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql-cnpg")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "openid")

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	env := []*k8s.EnvVar{
		{Name: jsii.String("PUID"), Value: jsii.String("1000")},
		{Name: jsii.String("PGID"), Value: jsii.String("1000")},
		{Name: jsii.String("TZ"), Value: jsii.String("Etc/UTC")},
		{
			Name: jsii.String("FRESHRSS_DB_HOST"),
			ValueFrom: &k8s.EnvVarSource{SecretKeyRef: &k8s.SecretKeySelector{
				Name: jsii.String("postgresql-cnpg"),
				Key:  jsii.String("host"),
			}},
		},
		{
			Name: jsii.String("FRESHRSS_DB_USER"),
			ValueFrom: &k8s.EnvVarSource{SecretKeyRef: &k8s.SecretKeySelector{
				Name: jsii.String("postgresql-cnpg"),
				Key:  jsii.String("username"),
			}},
		},
		{
			Name: jsii.String("FRESHRSS_DB_PASSWORD"),
			ValueFrom: &k8s.EnvVarSource{SecretKeyRef: &k8s.SecretKeySelector{
				Name: jsii.String("postgresql-cnpg"),
				Key:  jsii.String("password"),
			}},
		},
		{
			Name:  jsii.String("FRESHRSS_DB_NAME"),
			Value: jsii.String("freshrss"),
		},
	}

	// OIDC is done by Apache's mod_auth_openidc, which guards the web UI
	// (/i/) only: /api/ clients keep using their API password. Kept out of
	// env, which the updater CronJob shares and has no use for.
	oidcEnv := []*k8s.EnvVar{
		{Name: jsii.String("OIDC_ENABLED"), Value: jsii.String("1")},
		{
			Name:  jsii.String("OIDC_PROVIDER_METADATA_URL"),
			Value: jsii.String("https://auth.services.mkz.me/application/o/freshrss/.well-known/openid-configuration"),
		},
		{
			Name: jsii.String("OIDC_CLIENT_ID"),
			ValueFrom: &k8s.EnvVarSource{SecretKeyRef: &k8s.SecretKeySelector{
				Name: jsii.String("openid"),
				Key:  jsii.String("client_id"),
			}},
		},
		{
			Name: jsii.String("OIDC_CLIENT_SECRET"),
			ValueFrom: &k8s.EnvVarSource{SecretKeyRef: &k8s.SecretKeySelector{
				Name: jsii.String("openid"),
				Key:  jsii.String("client_secret"),
			}},
		},
		{
			Name: jsii.String("OIDC_CLIENT_CRYPTO_KEY"),
			ValueFrom: &k8s.EnvVarSource{SecretKeyRef: &k8s.SecretKeySelector{
				Name: jsii.String("openid"),
				Key:  jsii.String("crypto_key"),
			}},
		},
		{Name: jsii.String("OIDC_SCOPES"), Value: jsii.String("openid email profile")},
		// Traefik terminates TLS; without these the return URL is built as http://<pod>.
		{Name: jsii.String("OIDC_X_FORWARDED_HEADERS"), Value: jsii.String("X-Forwarded-Host X-Forwarded-Port X-Forwarded-Proto")},
		// The image default (300s) expires an idle tab's session within minutes; the
		// session stays capped by OIDC_SESSION_MAX_DURATION (image default: 27200s).
		{Name: jsii.String("OIDC_SESSION_INACTIVITY_TIMEOUT"), Value: jsii.String("28800")},
	}

	stsName, svcName := kubehelpers.NewStatefulSet(chart.Cdk8sChart, kubehelpers.StatefulSetConfig{
		Namespace: namespace,
		AppName:   appName,
		AppImage:  appImage,
		AppPort:   appPort,
		Labels:    labels,
		Env:       slices.Concat(env, oidcEnv),
		Storages: []kubehelpers.StatefulSetVolume{
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
	})

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
									jsii.String(stsName + "-0"),
								},
							},
						},
					},
				},
			},
		},
	}

	k8s.NewKubeCronJob(
		chart.Cdk8sChart,
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
										Env:  &env,
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
											ClaimName: jsii.String("data-" + stsName + "-0"),
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
		builder.Context,
		chart.Cdk8sChart,
		labels,
		appName,
		appPort,
		appIngress,
		svcName,
		map[string]string{},
	)

	return chart
}

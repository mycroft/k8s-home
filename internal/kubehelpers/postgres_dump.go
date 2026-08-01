package kubehelpers

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// PostgresDumpCronJobConfig describes a pg_dump CronJob and its backup PVC.
type PostgresDumpCronJobConfig struct {
	Namespace   string
	Name        string
	Image       string
	Schedule    string
	ClaimName   string
	StorageSize string
	Host        string
	Port        string
	SecretName  string
	Retention   int
}

// NewPostgresDumpCronJob creates daily logical dumps of every database in a
// PostgreSQL cluster. The job remains alive for eight hours so Velero's Kopia
// node agent can back up the mounted PVC after the dump has completed.
func NewPostgresDumpCronJob(chart cdk8s.Chart, cfg PostgresDumpCronJobConfig) {
	k8s.NewKubePersistentVolumeClaim(
		chart,
		jsii.String(cfg.ClaimName),
		&k8s.KubePersistentVolumeClaimProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String(cfg.ClaimName),
				Namespace: jsii.String(cfg.Namespace),
			},
			Spec: &k8s.PersistentVolumeClaimSpec{
				AccessModes:      &[]*string{jsii.String("ReadWriteOnce")},
				StorageClassName: jsii.String("longhorn-crypto-global"),
				Resources: &k8s.VolumeResourceRequirements{
					Requests: &map[string]k8s.Quantity{
						"storage": k8s.Quantity_FromString(jsii.String(cfg.StorageSize)),
					},
				},
			},
		},
	)

	command := fmt.Sprintf(`set -eu
umask 077
timestamp="$(date -u +%%Y%%m%%dT%%H%%M%%SZ)"
temporary="/backup/.${timestamp}.tmp"
destination="/backup/${timestamp}"
trap 'rm -rf "${temporary}"' EXIT
mkdir "${temporary}"
pg_dumpall --globals-only --file="${temporary}/globals.sql"
psql --dbname=postgres --tuples-only --no-align --command="SELECT datname FROM pg_database WHERE datallowconn AND NOT datistemplate ORDER BY datname" > "${temporary}/databases"
while IFS= read -r database; do
  pg_dump --format=custom --dbname="${database}" --file="${temporary}/${database}.dump"
done < "${temporary}/databases"
rm "${temporary}/databases"
mv "${temporary}" "${destination}"
trap - EXIT
ln -sfn "${timestamp}" /backup/latest
find /backup -mindepth 1 -maxdepth 1 -type d -mtime +%d -exec rm -rf {} +
# Keep the PVC mounted while Velero's scheduled filesystem backup runs.
sleep 8h`, cfg.Retention)

	k8s.NewKubeCronJob(
		chart,
		jsii.String(cfg.Name),
		&k8s.KubeCronJobProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String(cfg.Name),
				Namespace: jsii.String(cfg.Namespace),
			},
			Spec: &k8s.CronJobSpec{
				ConcurrencyPolicy:          jsii.String("Forbid"),
				FailedJobsHistoryLimit:     jsii.Number(1),
				Schedule:                   jsii.String(cfg.Schedule),
				SuccessfulJobsHistoryLimit: jsii.Number(1),
				JobTemplate: &k8s.JobTemplateSpec{
					Spec: &k8s.JobSpec{
						BackoffLimit: jsii.Number(1),
						Template: &k8s.PodTemplateSpec{
							Metadata: &k8s.ObjectMeta{
								Annotations: &map[string]*string{
									"backup.velero.io/backup-volumes": jsii.String("backup"),
								},
								Labels: &map[string]*string{
									"app.kubernetes.io/name": jsii.String(cfg.Name),
								},
							},
							Spec: &k8s.PodSpec{
								SecurityContext: &k8s.PodSecurityContext{
									FsGroup: jsii.Number(26),
								},
								Containers: &[]*k8s.Container{
									{
										Name:    jsii.String("postgres-dump"),
										Image:   jsii.String(cfg.Image),
										Command: &[]*string{jsii.String("/bin/sh"), jsii.String("-c"), jsii.String(command)},
										Env: &[]*k8s.EnvVar{
											{Name: jsii.String("PGHOST"), Value: jsii.String(cfg.Host)},
											{Name: jsii.String("PGPORT"), Value: jsii.String(cfg.Port)},
											{
												Name: jsii.String("PGUSER"),
												ValueFrom: &k8s.EnvVarSource{SecretKeyRef: &k8s.SecretKeySelector{
													Name: jsii.String(cfg.SecretName), Key: jsii.String("username"),
												}},
											},
											{
												Name: jsii.String("PGPASSWORD"),
												ValueFrom: &k8s.EnvVarSource{SecretKeyRef: &k8s.SecretKeySelector{
													Name: jsii.String(cfg.SecretName), Key: jsii.String("password"),
												}},
											},
										},
										VolumeMounts: &[]*k8s.VolumeMount{{Name: jsii.String("backup"), MountPath: jsii.String("/backup")}},
									},
								},
								RestartPolicy: jsii.String("Never"),
								Volumes: &[]*k8s.Volume{{
									Name:                  jsii.String("backup"),
									PersistentVolumeClaim: &k8s.PersistentVolumeClaimVolumeSource{ClaimName: jsii.String(cfg.ClaimName)},
								}},
							},
						},
					},
				},
			},
		},
	)
}

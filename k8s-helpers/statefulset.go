package k8s_helpers

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type StatefulSetVolume struct {
	Name        string
	MountPath   string
	StorageSize string
}

// NewStatefulSet creates a new statefulset and returns its name
func NewStatefulSet(
	chart cdk8s.Chart,
	namespace, appName, appImage string,
	appPort int,
	labels map[string]*string,
	env []*k8s.EnvVar,
	storages []StatefulSetVolume,
) string {
	svc := k8s.NewKubeService(
		chart,
		jsii.String("service"),
		&k8s.KubeServiceProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
				Labels:    &labels,
			},
			Spec: &k8s.ServiceSpec{
				Selector: &labels,
				Ports: &[]*k8s.ServicePort{
					{
						Port: jsii.Number(float64(appPort)),
						Name: jsii.String("http"),
					},
				},
			},
		},
	)

	mounts := []*k8s.VolumeMount{}
	pvcspecs := []*k8s.KubePersistentVolumeClaimProps{}
	for _, storage := range storages {
		mounts = append(mounts, &k8s.VolumeMount{
			MountPath: jsii.String(storage.MountPath),
			Name:      jsii.String(storage.Name),
		})
		pvcspecs = append(pvcspecs, &k8s.KubePersistentVolumeClaimProps{
			Metadata: &k8s.ObjectMeta{
				Labels: &labels,
				Name:   jsii.String(storage.Name),
			},
			Spec: &k8s.PersistentVolumeClaimSpec{
				AccessModes: &[]*string{
					jsii.String("ReadWriteOnce"),
				},
				StorageClassName: jsii.String("longhorn-crypto-global"),
				Resources: &k8s.ResourceRequirements{
					Requests: &map[string]k8s.Quantity{
						"storage": k8s.Quantity_FromString(jsii.String(storage.StorageSize)),
					},
				},
			},
		})
	}

	sts := k8s.NewKubeStatefulSet(
		chart,
		jsii.String("statefulset"),
		&k8s.KubeStatefulSetProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.StatefulSetSpec{
				Selector: &k8s.LabelSelector{
					MatchLabels: &labels,
				},
				ServiceName: svc.Name(),
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels: &labels,
					},
					Spec: &k8s.PodSpec{
						Containers: &[]*k8s.Container{
							{
								Name:         jsii.String(appName),
								Image:        jsii.String(appImage),
								Env:          &env,
								VolumeMounts: &mounts,
							},
						},
					},
				},
				VolumeClaimTemplates: &pvcspecs,
			},
		},
	)

	return *sts.Name()
}

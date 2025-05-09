package kubehelpers

import (
	"fmt"
	"strings"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type StatefulSetVolume struct {
	Name        string
	MountPath   string
	StorageSize string
}

type SecretMount struct {
	Name      string
	MountPath string
}

// NewStatefulSet creates a new statefulset and returns its name and its service name
func NewStatefulSet(
	chart cdk8s.Chart,
	namespace, appName, appImage string,
	appPort uint,
	labels map[string]*string,
	env []*k8s.EnvVar,
	commands []string,
	configMapMounts []ConfigMapMount,
	storages []StatefulSetVolume,
) (string, string) {
	return NewStatefulSetWithSecrets(chart, namespace, appName, appImage, appPort, labels, env, commands, configMapMounts, []SecretMount{}, storages)
}

// NewStatefulSet creates a new statefulset and returns its name and its service name
func NewStatefulSetWithSecrets(
	chart cdk8s.Chart,
	namespace, appName, appImage string,
	appPort uint,
	labels map[string]*string,
	env []*k8s.EnvVar,
	commands []string,
	configMapMounts []ConfigMapMount,
	secretMapMounts []SecretMount,
	storages []StatefulSetVolume,
) (string, string) {
	// Warning: Changing statefulSet object names will rename PVCs
	serviceObjectName := fmt.Sprintf("%s-svc", appName)
	statefulSetObjectName := fmt.Sprintf("%s-sts", appName)

	svc := k8s.NewKubeService(
		chart,
		jsii.String(serviceObjectName),
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

	volumeMounts := []*k8s.VolumeMount{}
	pvcspecs := []*k8s.KubePersistentVolumeClaimProps{}
	for _, storage := range storages {
		volumeMounts = append(volumeMounts, &k8s.VolumeMount{
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
				Resources: &k8s.VolumeResourceRequirements{
					Requests: &map[string]k8s.Quantity{
						"storage": k8s.Quantity_FromString(jsii.String(storage.StorageSize)),
					},
				},
			},
		})
	}

	volumes := []*k8s.Volume{}

	for _, secret := range secretMapMounts {
		volumes = append(volumes, &k8s.Volume{
			Name: jsii.String(secret.Name),
			Secret: &k8s.SecretVolumeSource{
				SecretName: jsii.String(secret.Name),
			},
		})

		volumeMounts = append(volumeMounts, &k8s.VolumeMount{
			Name:      jsii.String(secret.Name),
			MountPath: jsii.String(secret.MountPath),
		})
	}

	for _, v := range configMapMounts {
		volumes = append(volumes, &k8s.Volume{
			Name: jsii.String(v.Name),
			ConfigMap: &k8s.ConfigMapVolumeSource{
				Name: v.ConfigMap.Name(),
			},
		})

		volumeMounts = append(volumeMounts, &k8s.VolumeMount{
			Name:      jsii.String(v.Name),
			MountPath: jsii.String(v.MountPath),
		})
	}

	container := k8s.Container{
		Name:         jsii.String(appName),
		Image:        jsii.String(appImage),
		Env:          &env,
		VolumeMounts: &volumeMounts,
	}

	// If only one command...
	if len(commands) == 1 {
		commandsElmts := strings.Split(commands[0], " ")
		command := []*string{}
		for _, el := range commandsElmts {
			command = append(command, jsii.String(el))
		}
		container.Command = &command
	} else if len(commands) > 0 { // or multiple...
		container.Command = &[]*string{
			jsii.String("/bin/sh"),
			jsii.String("-c"),
			jsii.String(strings.Join(commands, " && ")),
		}
	}

	sts := k8s.NewKubeStatefulSet(
		chart,
		jsii.String(statefulSetObjectName),
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
							&container,
						},
						Volumes: &volumes,
					},
				},
				VolumeClaimTemplates: &pvcspecs,
			},
		},
	)

	return *sts.Name(), *svc.Name()
}

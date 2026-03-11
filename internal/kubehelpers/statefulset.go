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

type StatefulSetConfig struct {
	Namespace       string
	AppName         string
	AppImage        string
	AppPort         uint
	Labels          map[string]*string
	Env             []*k8s.EnvVar
	Commands        []string
	ConfigMapMounts []ConfigMapMount
	SecretMounts    []SecretMount
	Storages        []StatefulSetVolume
	FsGroup         int
}

// NewStatefulSet creates a new statefulset and returns its name and its service name
func NewStatefulSet(chart cdk8s.Chart, cfg StatefulSetConfig) (string, string) {
	// Warning: Changing statefulSet object names will rename PVCs
	serviceObjectName := fmt.Sprintf("%s-svc", cfg.AppName)
	statefulSetObjectName := fmt.Sprintf("%s-sts", cfg.AppName)

	svc := k8s.NewKubeService(
		chart,
		jsii.String(serviceObjectName),
		&k8s.KubeServiceProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(cfg.Namespace),
				Labels:    &cfg.Labels,
			},
			Spec: &k8s.ServiceSpec{
				Selector: &cfg.Labels,
				Ports: &[]*k8s.ServicePort{
					{
						Port: jsii.Number(float64(cfg.AppPort)),
						Name: jsii.String("http"),
					},
				},
			},
		},
	)

	volumeMounts := []*k8s.VolumeMount{}
	pvcspecs := []*k8s.KubePersistentVolumeClaimProps{}
	for _, storage := range cfg.Storages {
		volumeMounts = append(volumeMounts, &k8s.VolumeMount{
			MountPath: jsii.String(storage.MountPath),
			Name:      jsii.String(storage.Name),
		})
		pvcspecs = append(pvcspecs, &k8s.KubePersistentVolumeClaimProps{
			Metadata: &k8s.ObjectMeta{
				Labels: &cfg.Labels,
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

	for _, secret := range cfg.SecretMounts {
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

	for _, v := range cfg.ConfigMapMounts {
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
		Name:         jsii.String(cfg.AppName),
		Image:        jsii.String(cfg.AppImage),
		Env:          &cfg.Env,
		VolumeMounts: &volumeMounts,
	}

	// If only one command...
	if len(cfg.Commands) == 1 {
		commandsElmts := strings.Split(cfg.Commands[0], " ")
		command := []*string{}
		for _, el := range commandsElmts {
			command = append(command, jsii.String(el))
		}
		container.Command = &command
	} else if len(cfg.Commands) > 0 { // or multiple...
		container.Command = &[]*string{
			jsii.String("/bin/sh"),
			jsii.String("-c"),
			jsii.String(strings.Join(cfg.Commands, " && ")),
		}
	}

	sts := k8s.NewKubeStatefulSet(
		chart,
		jsii.String(statefulSetObjectName),
		&k8s.KubeStatefulSetProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(cfg.Namespace),
			},
			Spec: &k8s.StatefulSetSpec{
				Selector: &k8s.LabelSelector{
					MatchLabels: &cfg.Labels,
				},
				ServiceName: svc.Name(),
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels: &cfg.Labels,
					},
					Spec: &k8s.PodSpec{
						Containers: &[]*k8s.Container{
							&container,
						},
						SecurityContext: &k8s.PodSecurityContext{
							FsGroup: jsii.Number(cfg.FsGroup),
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

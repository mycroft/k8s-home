package kubehelpers

import (
	"fmt"
	"strings"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// StatefulSetVolume describes a persistent volume to be created and mounted into the container.
// Each entry produces a PersistentVolumeClaim (using the longhorn-crypto-global storage class)
// and a corresponding VolumeMount inside the StatefulSet pod.
type StatefulSetVolume struct {
	Name        string // PVC and volume name; also used as the VolumeMount name
	MountPath   string // Absolute path inside the container where the volume is mounted
	StorageSize string // Kubernetes quantity string, e.g. "10Gi"
}

// SecretMount describes a Kubernetes Secret to be projected as a volume into the container.
// The secret must already exist in the namespace (e.g. created by ExternalSecret).
type SecretMount struct {
	Name      string // Name of the existing Kubernetes Secret and the resulting volume
	MountPath string // Directory inside the container where the secret keys are mounted as files
}

// StatefulSetConfig is the input to NewStatefulSet. It describes a single-replica StatefulSet
// together with its headless Service, persistent storage, and optional environment/mounts.
type StatefulSetConfig struct {
	Namespace string // Kubernetes namespace for all generated objects
	AppName   string // Base name used to derive object names (e.g. "<AppName>-sts", "<AppName>-svc")

	AppImage string // Container image reference, including tag (sourced from versions.yaml via RegisterContainerImage)
	AppPort  uint   // Port the container listens on; exposed by the Service and used by Ingress helpers

	// Labels are applied to the Pod template, the Service selector, and PVC metadata.
	// Convention: {"app.kubernetes.io/name": appName}
	Labels map[string]*string

	// Env is a list of environment variables injected into the container.
	Env []*k8s.EnvVar

	// Commands overrides the container's default entrypoint.
	// - Zero entries: no override (image ENTRYPOINT is used).
	// - One entry: the string is split on spaces and passed as container.command.
	// - Two or more entries: executed as `/bin/sh -c "cmd1 && cmd2 && ..."`.
	Commands []string

	// ConfigMapMounts projects ConfigMap data as files into the container.
	ConfigMapMounts []ConfigMapMount

	// SecretMounts projects Kubernetes Secrets as files into the container.
	SecretMounts []SecretMount

	// Storages defines the persistent volumes attached to the StatefulSet.
	// Each entry creates a PVC (ReadWriteOnce, longhorn-crypto-global) and mounts it.
	Storages []StatefulSetVolume

	// FsGroup sets the pod-level fsGroup in the SecurityContext so that mounted
	// volumes are owned by the specified GID. Set to 0 to leave unset.
	FsGroup int
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

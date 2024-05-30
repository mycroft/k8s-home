package kubehelpers

import (
	"strings"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type ConfigMapMount struct {
	Name      string
	ConfigMap k8s.KubeConfigMap
	MountPath string
}

func NewAppDeployment(
	chart cdk8s.Chart,
	namespace, appName, appImage string,
	labels map[string]*string,
	env []*k8s.EnvVar,
	commands []string,
	configMapMounts []ConfigMapMount,
) {

	volumes := []*k8s.Volume{}
	volumeMounts := []*k8s.VolumeMount{}

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

	if len(commands) == 1 { // if one command...
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

	k8s.NewKubeDeployment(
		chart,
		jsii.String("deploy"),
		&k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.DeploymentSpec{
				Selector: &k8s.LabelSelector{
					MatchLabels: &labels,
				},
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
			},
		},
	)
}

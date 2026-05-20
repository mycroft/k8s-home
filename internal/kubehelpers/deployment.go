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

type AppDeploymentOption struct {
	Name            string
	ImagePullPolicy string
}

func NewAppDeployment(
	chart cdk8s.Chart,
	namespace, appName, appImage string,
	labels map[string]*string,
	env []*k8s.EnvVar,
	commands []string,
	configMapMounts []ConfigMapMount,
	opts ...AppDeploymentOption,
) {
	imagePullPolicy := ""

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

	metadatas := k8s.ObjectMeta{
		Namespace: jsii.String(namespace),
	}

	for _, opt := range opts {
		if opt.Name != "" {
			metadatas.Name = jsii.String(opt.Name)
		}

		if opt.ImagePullPolicy != "" {
			imagePullPolicy = opt.ImagePullPolicy
		}
	}

	container := k8s.Container{
		Name:  jsii.String(appName),
		Image: jsii.String(appImage),
	}

	if imagePullPolicy != "" {
		container.ImagePullPolicy = jsii.String(imagePullPolicy)
	}

	if len(env) > 0 {
		container.Env = &env
	}

	if len(volumeMounts) > 0 {
		container.VolumeMounts = &volumeMounts
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

	podTemplateSpec := k8s.PodTemplateSpec{
		Metadata: &k8s.ObjectMeta{
			Labels: &labels,
		},
		Spec: &k8s.PodSpec{
			Containers: &[]*k8s.Container{
				&container,
			},
		},
	}

	if len(volumes) > 0 {
		podTemplateSpec.Spec.Volumes = &volumes
	}

	k8s.NewKubeDeployment(
		chart,
		jsii.String("deploy"),
		&k8s.KubeDeploymentProps{
			Metadata: &metadatas,
			Spec: &k8s.DeploymentSpec{
				Selector: &k8s.LabelSelector{
					MatchLabels: &labels,
				},
				Template: &podTemplateSpec,
			},
		},
	)
}

// Deployment contains a deployment configuration
type Deployment struct {
	Name            string
	Image           string
	ImagePullPolicy string
	// Labels to apply
	Labels map[string]*string
	// Environement to set in deployment's pods
	Env        []*k8s.EnvVar
	Commands   []string
	ConfigMaps []ConfigMapMount
}

func (chart *Chart) NewDeployment(deployment *Deployment) {
	if chart.Namespace == "" {
		panic("namespace was not defined")
	}

	NewAppDeployment(
		chart.Cdk8sChart,
		chart.Namespace,
		deployment.Name,
		chart.Builder.RegisterContainerImage(deployment.Image),
		deployment.Labels,
		deployment.Env,
		deployment.Commands,
		deployment.ConfigMaps,
		AppDeploymentOption{
			Name:            deployment.Name,
			ImagePullPolicy: deployment.ImagePullPolicy,
		},
	)
}

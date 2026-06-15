package observability

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

const (
	smokepingProberConfigDir  = "/etc/smokeping_prober"
	smokepingProberConfigFile = "smokeping_prober.yml"
)

func NewSmokepingProberChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	app := "smokeping-prober"
	namespace := app
	image := builder.RegisterContainerImage("quay.io/superq/smokeping-prober")

	chart := builder.NewChart(app)
	chart.NewNamespace(namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(app),
	}

	configMap := k8s.NewKubeConfigMap(
		chart.Cdk8sChart,
		jsii.String("config"),
		&k8s.KubeConfigMapProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
			},
			Data: &map[string]*string{
				smokepingProberConfigFile: jsii.String(`---
targets:
- hosts:
  - moonstone.lan.mkz.me
  - glitter.lan.mkz.me
  interval: 5s
`),
			},
		},
	)

	k8s.NewKubeDeployment(
		chart.Cdk8sChart,
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
						Annotations: &map[string]*string{
							"configMapHash": jsii.String(kubehelpers.ComputeConfigMapHash(configMap)),
						},
					},
					Spec: &k8s.PodSpec{
						Containers: &[]*k8s.Container{
							{
								Name:  jsii.String(app),
								Image: jsii.String(image),
								Args: &[]*string{
									jsii.String("--config.file=" + smokepingProberConfigDir + "/" + smokepingProberConfigFile),
								},
								Ports: &[]*k8s.ContainerPort{
									{
										Name:          jsii.String("http"),
										ContainerPort: jsii.Number(9374),
									},
								},
								SecurityContext: &k8s.SecurityContext{
									Capabilities: &k8s.Capabilities{
										Add: &[]*string{
											jsii.String("NET_RAW"),
										},
									},
								},
								VolumeMounts: &[]*k8s.VolumeMount{
									{
										Name:      jsii.String("config"),
										MountPath: jsii.String(smokepingProberConfigDir),
									},
								},
							},
						},
						Volumes: &[]*k8s.Volume{
							{
								Name: jsii.String("config"),
								ConfigMap: &k8s.ConfigMapVolumeSource{
									Name: configMap.Name(),
								},
							},
						},
					},
				},
			},
		},
	)

	kubehelpers.NewAppService(
		chart.Cdk8sChart,
		namespace,
		"svc",
		labels,
		"http",
		9374,
	)

	chart.NewServiceMonitor(&kubehelpers.ServiceMonitor{
		Labels: map[string]string{
			"app.kubernetes.io/name": app,
		},
	})

	return chart
}

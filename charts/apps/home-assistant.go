package apps

import (
	"os"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

// NewHomeAssistantChart runs Home Assistant on the host network: device
// discovery (mDNS, SSDP, DHCP) needs to see the LAN's multicast and broadcast
// traffic, which the pod network does not carry. It therefore also listens on
// <node>:8123, besides the ingress.
//
// kubehelpers.NewStatefulSet covers none of hostNetwork, init containers or a
// file mounted inside the data volume, hence the hand-written StatefulSet.
func NewHomeAssistantChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "home-assistant"
	appName := namespace
	appPort := uint(8123)
	appIngress := "home-assistant.iop.cx"
	image := builder.RegisterContainerImage("ghcr.io/home-assistant/home-assistant")

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	configuration, err := os.ReadFile("configs/home-assistant/configuration.yaml")
	if err != nil {
		panic(err)
	}

	configMap := k8s.NewKubeConfigMap(
		chart.Cdk8sChart,
		jsii.String("config"),
		&k8s.KubeConfigMapProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
			},
			Data: &map[string]*string{
				"configuration.yaml": jsii.String(string(configuration)),
			},
		},
	)

	svc := kubehelpers.NewAppService(
		chart.Cdk8sChart,
		namespace,
		"svc",
		labels,
		"http",
		appPort,
	)

	dataMount := &k8s.VolumeMount{
		Name:      jsii.String("data"),
		MountPath: jsii.String("/config"),
	}

	k8s.NewKubeStatefulSet(
		chart.Cdk8sChart,
		jsii.String("sts"),
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
						// The ConfigMap is mounted with subPath, which the
						// kubelet never refreshes: roll the pod instead.
						Annotations: &map[string]*string{
							"configMapHash": jsii.String(kubehelpers.ComputeConfigMapHash(configMap)),
						},
					},
					Spec: &k8s.PodSpec{
						HostNetwork: jsii.Bool(true),
						// Keep resolving cluster names despite hostNetwork.
						DnsPolicy: jsii.String("ClusterFirstWithHostNet"),
						InitContainers: &[]*k8s.Container{
							{
								Name:  jsii.String("init-includes"),
								Image: jsii.String(image),
								Command: &[]*string{
									jsii.String("/bin/sh"),
									jsii.String("-c"),
									// Same contents as Home Assistant's own defaults.
									jsii.String("cd /config && { [ -e automations.yaml ] || echo '[]' > automations.yaml; } && touch scripts.yaml scenes.yaml"),
								},
								VolumeMounts: &[]*k8s.VolumeMount{dataMount},
							},
						},
						Containers: &[]*k8s.Container{
							{
								Name:  jsii.String(appName),
								Image: jsii.String(image),
								Ports: &[]*k8s.ContainerPort{
									{
										Name:          jsii.String("http"),
										ContainerPort: jsii.Number(float64(appPort)),
									},
								},
								VolumeMounts: &[]*k8s.VolumeMount{
									dataMount,
									{
										Name:      jsii.String("config"),
										MountPath: jsii.String("/config/configuration.yaml"),
										SubPath:   jsii.String("configuration.yaml"),
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
				VolumeClaimTemplates: &[]*k8s.KubePersistentVolumeClaimProps{
					{
						Metadata: &k8s.ObjectMeta{
							Name:   jsii.String("data"),
							Labels: &labels,
						},
						Spec: &k8s.PersistentVolumeClaimSpec{
							AccessModes: &[]*string{
								jsii.String("ReadWriteOnce"),
							},
							StorageClassName: jsii.String("longhorn-crypto-global"),
							Resources: &k8s.VolumeResourceRequirements{
								Requests: &map[string]k8s.Quantity{
									"storage": k8s.Quantity_FromString(jsii.String("10Gi")),
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
		*svc.Name(),
		map[string]string{},
	)

	return chart
}

package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

const (
	happyUrlsImage = "git.mkz.me/mycroft/happy-urls:latest"
)

func NewHappyUrlsChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "happy-urls"
	appName := namespace
	appPort := uint(3000)
	ingressHost := fmt.Sprintf("%s.services.mkz.me", appName)

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	redisLabels := map[string]*string{
		"app.kubernetes.io/name":      jsii.String(appName),
		"app.kubernetes.io/component": jsii.String("redis"),
	}

	_, redisSvcName := kubehelpers.NewStatefulSet(
		chart.Cdk8sChart,
		namespace,
		"redis",
		builder.RegisterContainerImage("redis"),
		6379,
		redisLabels,
		[]*k8s.EnvVar{},
		[]string{
			"redis-server --save 60 1 --loglevel warning",
		},
		[]kubehelpers.ConfigMapMount{},
		[]kubehelpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/data",
				StorageSize: "1Gi",
			},
		},
	)

	labels := map[string]*string{
		"app.kubernetes.io/name":      jsii.String(appName),
		"app.kubernetes.io/component": jsii.String(appName),
	}

	envVars := []*k8s.EnvVar{
		{
			Name:  jsii.String("REDIS_URL"),
			Value: jsii.String(fmt.Sprintf("redis://%s:6379", redisSvcName)),
		},
	}

	k8s.NewKubeDeployment(
		chart.Cdk8sChart,
		jsii.String("deploy"),
		&k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(appName),
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
							{
								Env:             &envVars,
								ImagePullPolicy: jsii.String("Always"),
								Name:            jsii.String(appName),
								Image:           jsii.String(happyUrlsImage),
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
		ingressHost,
		"",
		map[string]string{},
	)

	return chart
}

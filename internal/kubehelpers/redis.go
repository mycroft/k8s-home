package kubehelpers

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/jsii-runtime-go"
)

// NewRedisStatefulset creates a statefulset of redis in given chart and returns its name and service name
func (chart *Chart) NewRedisStatefulset(namespace string) (string, string) {
	imageName := chart.Builder.RegisterContainerImage("redis")

	redisLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("redis"),
	}

	sst, serviceName := NewStatefulSet(
		chart.Cdk8sChart,
		namespace,
		"redis",
		imageName,
		6379,
		redisLabels,
		[]*k8s.EnvVar{},
		[]string{
			"redis-server --save 60 1 --loglevel warning",
		},
		[]ConfigMapMount{},
		[]StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/data",
				StorageSize: "1Gi",
			},
		},
	)

	return sst, serviceName

}

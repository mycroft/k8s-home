package kubehelpers

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// NewRedisStatefulset creates a statefulset of redis in given chart and returns its name and service name
func NewRedisStatefulset(chart cdk8s.Chart, namespace string) (string, string) {
	redisLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("redis"),
	}

	sst, serviceName := NewStatefulSet(
		chart,
		namespace,
		"redis",
		RegisterDockerImage("redis"),
		6379,
		redisLabels,
		[]*k8s.EnvVar{},
		[]string{
			"redis-server --save 60 1 --loglevel warning",
		},
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

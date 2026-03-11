package kubehelpers

import (
	"github.com/aws/jsii-runtime-go"
)

// NewRedisStatefulset creates a statefulset of redis in given chart and returns its name and service name
func (chart *Chart) NewRedisStatefulset(namespace string) (string, string) {
	imageName := chart.Builder.RegisterContainerImage("redis")

	redisLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("redis"),
	}

	sst, serviceName := NewStatefulSet(chart.Cdk8sChart, StatefulSetConfig{
		Namespace: namespace,
		AppName:   "redis",
		AppImage:  imageName,
		AppPort:   6379,
		Labels:    redisLabels,
		Commands: []string{
			"redis-server --save 60 1 --loglevel warning",
		},
		Storages: []StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/data",
				StorageSize: "1Gi",
			},
		},
	})

	return sst, serviceName

}

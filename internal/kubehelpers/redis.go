package kubehelpers

import (
	"fmt"
	"maps"

	"github.com/aws/jsii-runtime-go"
)

// RedisOpts configures optional parameters for a Redis StatefulSet.
// Zero values fall back to the defaults used by NewRedisStatefulset.
type RedisOpts struct {
	AppPort     uint
	Labels      map[string]*string
	StorageSize string
	// LogLevel is a Redis log verbosity level: debug, verbose, notice, warning.
	LogLevel string
}

// NewRedisStatefulset provisions a single-replica Redis StatefulSet in the given namespace.
// Redis runs with RDB persistence (snapshot every 60 s if at least 1 key changed, via --save 60 1)
// and warning-level logging. Storage is fixed at 1 Gi on the cluster encrypted storage class.
// Returns the StatefulSet name and the ClusterIP service name.
func (chart *Chart) NewRedisStatefulset(namespace string) (string, string) {
	return chart.NewRedisStatefulsetWithOpts(namespace, RedisOpts{})
}

// NewRedisStatefulsetWithOpts is like NewRedisStatefulset but accepts a RedisOpts
// to override defaults: port (6379), extra labels, storage size (1Gi), and log level (warning).
// Returns the StatefulSet name and the ClusterIP service name.
func (chart *Chart) NewRedisStatefulsetWithOpts(namespace string, opts RedisOpts) (string, string) {
	imageName := chart.Builder.RegisterContainerImage("redis")

	appPort := opts.AppPort
	if appPort == 0 {
		appPort = 6379
	}

	storageSize := opts.StorageSize
	if storageSize == "" {
		storageSize = "1Gi"
	}

	logLevel := opts.LogLevel
	if logLevel == "" {
		logLevel = "warning"
	}

	redisLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("redis"),
	}
	maps.Copy(redisLabels, opts.Labels)

	sst, serviceName := NewStatefulSet(chart.Cdk8sChart, StatefulSetConfig{
		Namespace: namespace,
		AppName:   "redis",
		AppImage:  imageName,
		AppPort:   appPort,
		Labels:    redisLabels,
		Commands: []string{
			fmt.Sprintf("redis-server --save 60 1 --loglevel %s", logLevel),
		},
		Storages: []StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/data",
				StorageSize: storageSize,
			},
		},
	})

	return sst, serviceName
}

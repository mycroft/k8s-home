package kubehelpers_test

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// ExampleNewStatefulSet demonstrates how to deploy a stateful application with
// persistent storage, environment variables, a startup command, and a non-root
// filesystem group. The returned names are typically passed to NewAppIngress.
func ExampleNewStatefulSet() {
	app := cdk8s.NewApp(nil)
	chart := cdk8s.NewChart(app, jsii.String("myapp"), &cdk8s.ChartProps{})

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String("myapp"),
	}

	stsName, svcName := kubehelpers.NewStatefulSet(chart, kubehelpers.StatefulSetConfig{
		Namespace: "myapp",
		AppName:   "myapp",
		AppImage:  "myorg/myapp:1.0.0",
		AppPort:   8080,
		Labels:    labels,

		// Environment variables injected into the container.
		Env: []*k8s.EnvVar{
			{
				Name:  jsii.String("DATABASE_URL"),
				Value: jsii.String("postgres://postgres-instance.postgres:5432/myapp"),
			},
			{
				// Source a sensitive value from an existing Kubernetes Secret
				// (e.g. one created by CreateExternalSecret).
				Name: jsii.String("DATABASE_PASSWORD"),
				ValueFrom: &k8s.EnvVarSource{
					SecretKeyRef: &k8s.SecretKeySelector{
						Name: jsii.String("postgresql"),
						Key:  jsii.String("password"),
					},
				},
			},
		},

		// Single-string command: split on spaces and used as container.command.
		// Use a slice of multiple strings to join them with " && " under /bin/sh -c.
		Commands: []string{"/app/server --config /etc/myapp/config.yaml"},

		// Each entry provisions a PVC (ReadWriteOnce, longhorn-crypto-global)
		// and mounts it at the given path inside the container.
		Storages: []kubehelpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/var/lib/myapp/data",
				StorageSize: "10Gi",
			},
		},

		// Set the pod fsGroup so mounted volumes are group-owned by GID 1000.
		FsGroup: 1000,
	})

	// stsName and svcName are the synthesised Kubernetes object names.
	// Pass svcName to NewAppIngress to route external traffic to this StatefulSet.
	_ = stsName
	_ = svcName
}

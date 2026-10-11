package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

// NewZulipChart deploys Zulip (https://zulip.com), the team chat used to talk
// to agent-nexus, from its upstream Helm chart. Values are in
// configs/zulip.yaml.tmpl.
//
// The chart's Bitnami subcharts are all disabled (their images are the
// unmaintained bitnamilegacy ones), and Zulip's services come from here:
//
//   - PostgreSQL: the zulip database on the CNPG cluster
//     (charts/storage/cnpg-cluster.go), credentials from Vault at
//     secret/namespaces/zulip/postgresql;
//   - Redis: the shared Redis helper;
//   - RabbitMQ and memcached: single-replica StatefulSets below.
//
// Outgoing email goes through the mail.mkz.me relay, with credentials in
// Vault at secret/namespaces/zulip/mailrelay (username, password).
func NewZulipChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "zulip"
	repositoryName := "zulip"
	chartName := "zulip"
	releaseName := "zulip"
	appIngress := "zulip.services.mkz.me"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "mailrelay")
	// Django's SECRET_KEY, which the chart requires. Pinned once generated:
	// changing it logs everyone out.
	kubehelpers.CreateGeneratedPasswordSecret(chart.Cdk8sChart, namespace, "secret-key")

	_, redisServiceName := chart.NewRedisStatefulset(namespace)
	rabbitmqServiceName := newZulipRabbitMQ(chart, namespace)
	memcachedServiceName := newZulipMemcached(chart, namespace)

	chart.CreateHelmRepository(repositoryName, "oci://ghcr.io/zulip/helm-charts")

	configMaps := []kubehelpers.HelmReleaseConfigMap{
		chart.CreateHelmValuesConfig(
			namespace,
			releaseName,
			releaseName+".yaml",
			map[string]interface{}{
				"ExternalHost": appIngress,
				"Redis":        redisServiceName,
				"RabbitMQ":     rabbitmqServiceName,
				"Memcached":    memcachedServiceName,
			},
		),
	}

	// Covers the 10 minutes the values give the first start (startupProbe).
	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithConfigMaps(configMaps),
		kubehelpers.WithTimeout("15m0s"),
	)

	chart.NewIngress(&kubehelpers.Ingress{
		Name:        "zulip",
		Ingresses:   []string{appIngress},
		ServiceName: releaseName, // created by the helm chart
	})

	return chart
}

// newZulipRabbitMQ runs the broker Zulip queues its deferred work on: emails,
// push notifications, and the events sent to outgoing webhook bots. Its data
// is persistent so a restart does not drop pending events, which also means
// the zulip user is only created on the first start, from the generated
// password. Returns the service name.
func newZulipRabbitMQ(chart *kubehelpers.Chart, namespace string) string {
	secretName := "rabbitmq"
	kubehelpers.CreateGeneratedPasswordSecret(chart.Cdk8sChart, namespace, secretName)

	_, serviceName := kubehelpers.NewStatefulSet(chart.Cdk8sChart, kubehelpers.StatefulSetConfig{
		Namespace: namespace,
		AppName:   "rabbitmq",
		AppImage:  chart.Builder.RegisterContainerImage("rabbitmq"),
		AppPort:   5672,
		Labels: map[string]*string{
			"app.kubernetes.io/component": jsii.String("rabbitmq"),
		},
		Env: []*k8s.EnvVar{
			{
				Name:  jsii.String("RABBITMQ_DEFAULT_USER"),
				Value: jsii.String("zulip"),
			},
			{
				Name: jsii.String("RABBITMQ_DEFAULT_PASS"),
				ValueFrom: &k8s.EnvVarSource{
					SecretKeyRef: &k8s.SecretKeySelector{
						Name: jsii.String(secretName),
						Key:  jsii.String("password"),
					},
				},
			},
		},
		Storages: []kubehelpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/var/lib/rabbitmq",
				StorageSize: "1Gi",
			},
		},
	})

	return serviceName
}

// newZulipMemcached runs Zulip's cache. It has no storage and no SASL: Zulip
// only authenticates to memcached when given a password. A StatefulSet, like
// the Redis helper, for a single stable pod behind its own service. Returns
// the service name.
func newZulipMemcached(chart *kubehelpers.Chart, namespace string) string {
	_, serviceName := kubehelpers.NewStatefulSet(chart.Cdk8sChart, kubehelpers.StatefulSetConfig{
		Namespace: namespace,
		AppName:   "memcached",
		AppImage:  chart.Builder.RegisterContainerImage("memcached"),
		AppPort:   11211,
		Labels: map[string]*string{
			"app.kubernetes.io/component": jsii.String("memcached"),
		},
		Commands: []string{
			"memcached --memory-limit=256",
		},
	})

	return serviceName
}

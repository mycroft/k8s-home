package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/imports/servicemonitor_monitoringcoreoscom"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewYopassChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "yopass"
	appIngress := "yopass.services.mkz.me"
	appName := "yopass"
	appPort := 1337
	image := kubehelpers.RegisterDockerImage("jhaals/yopass")

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	_, redisServiceName := kubehelpers.NewRedisStatefulset(chart.Cdk8sChart, namespace)
	redisURL := fmt.Sprintf("redis://%s:6379", redisServiceName)

	yopassLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("yopass"),
	}

	kubehelpers.NewAppDeployment(
		chart.Cdk8sChart,
		namespace,
		appName,
		image,
		yopassLabels,
		[]*k8s.EnvVar{},
		[]string{
			fmt.Sprintf("/yopass-server --database redis --metrics-port 1338 --port 1337 --redis %s", redisURL),
		},
		[]kubehelpers.ConfigMapMount{},
	)

	kubehelpers.NewAppIngress(
		builder.Context,
		chart.Cdk8sChart,
		yopassLabels,
		appName,
		appPort,
		appIngress,
		"",
		map[string]string{},
	)

	kubehelpers.NewAppService(
		chart.Cdk8sChart,
		namespace,
		"metrics",
		yopassLabels,
		"metrics",
		1338,
	)

	servicemonitor_monitoringcoreoscom.NewServiceMonitor(
		chart.Cdk8sChart,
		jsii.String("sm"),
		&servicemonitor_monitoringcoreoscom.ServiceMonitorProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Labels: &map[string]*string{
					"release": jsii.String("prometheus"),
				},
			},
			Spec: &servicemonitor_monitoringcoreoscom.ServiceMonitorSpec{
				Selector: &servicemonitor_monitoringcoreoscom.ServiceMonitorSpecSelector{
					MatchLabels: &map[string]*string{},
				},
				Endpoints: &[]*servicemonitor_monitoringcoreoscom.ServiceMonitorSpecEndpoints{
					{
						Port: jsii.String("metrics"),
					},
				},
			},
		},
	)

	return chart
}

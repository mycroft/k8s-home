package apps

import (
	"context"
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/imports/servicemonitor_monitoringcoreoscom"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewYopassChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	namespace := "yopass"
	appIngress := "yopass.services.mkz.me"
	appName := "yopass"
	appPort := 1337
	image := kubehelpers.RegisterDockerImage("jhaals/yopass")

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	_, redisServiceName := kubehelpers.NewRedisStatefulset(chart, namespace)
	redisURL := fmt.Sprintf("redis://%s:6379", redisServiceName)

	yopassLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("yopass"),
	}

	kubehelpers.NewAppDeployment(
		chart,
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
		chart,
		yopassLabels,
		appName,
		appPort,
		appIngress,
		"",
		map[string]string{},
	)

	kubehelpers.NewAppService(
		chart,
		namespace,
		"metrics",
		yopassLabels,
		"metrics",
		1338,
	)

	servicemonitor_monitoringcoreoscom.NewServiceMonitor(
		chart,
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

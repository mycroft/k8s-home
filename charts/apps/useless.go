package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/imports/servicemonitor_monitoringcoreoscom"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewUselessChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	name := "useless"

	namespace := name
	appName := name
	appImage := builder.RegisterContainerImage("registry.mkz.me/mycroft/useless")
	appPort := uint(8080)
	appIngress := "useless.iop.cx"

	chart := builder.NewChart(name)
	chart.NewNamespace(namespace)

	labels := map[string]*string{
		"app.kubernetes.io/name":      jsii.String(appName),
		"app.kubernetes.io/component": jsii.String("api"),
	}

	redisOpts := kubehelpers.RedisOpts{
		AppPort: 6379,
	}
	_, redisServiceName := chart.NewRedisStatefulsetWithOpts(namespace, redisOpts)

	env := []*k8s.EnvVar{
		{
			Name:  jsii.String("REDIS_HOST"),
			Value: jsii.Sprintf("%s.%s", redisServiceName, namespace),
		},
		{
			Name:  jsii.String("REDIS_PORT"),
			Value: jsii.Sprintf("%d", redisOpts.AppPort),
		},
	}

	k8s.NewKubeDeployment(
		chart.Cdk8sChart,
		jsii.String("useless"),
		&k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
				Labels:    &labels,
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
								Env:             &env,
								Name:            jsii.String(appName),
								Image:           jsii.String(appImage),
								ImagePullPolicy: jsii.String("Always"),
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
		"",
		map[string]string{},
	)

	servicemonitor_monitoringcoreoscom.NewServiceMonitor(
		chart.Cdk8sChart,
		jsii.String("service-monitor"),
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
						Path: jsii.String("/metrics"),
						Port: jsii.String("http"),
					},
				},
			},
		},
	)

	return chart
}

package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewUselessChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	name := "useless"

	namespace := name
	appName := name
	appImage := builder.RegisterContainerImage("registry.mkz.me/mycroft/useless")
	appPort := 8080
	appIngress := "useless.iop.cx"
	redisPort := 6379

	chart := builder.NewChart(name)
	chart.NewNamespace(namespace)

	labels := map[string]string{
		"app.kubernetes.io/name":      appName,
		"app.kubernetes.io/component": "api",
	}

	_, redisServiceName := chart.NewRedisStatefulsetWithOpts(namespace, kubehelpers.RedisOpts{
		AppPort: uint(redisPort),
	})

	env := []*k8s.EnvVar{
		{
			Name:  jsii.String("REDIS_HOST"),
			Value: jsii.Sprintf("%s.%s", redisServiceName, namespace),
		},
		{
			Name:  jsii.String("REDIS_PORT"),
			Value: jsii.Sprintf("%d", redisPort),
		},
	}

	chart.NewDeployment(&kubehelpers.Deployment{
		Name:            name,
		Labels:          labels,
		Env:             env,
		Image:           appImage,
		ImagePullPolicy: "Always",
	})

	chart.NewIngress(&kubehelpers.Ingress{
		Name:   name,
		Labels: labels,
		Ingresses: []string{
			appIngress,
		},
		Port: uint(appPort),
	})

	chart.NewServiceMonitor(&kubehelpers.ServiceMonitor{
		Labels: labels,
	})

	return chart
}

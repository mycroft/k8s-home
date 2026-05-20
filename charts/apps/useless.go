package apps

import (
	"fmt"

	kube "git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewUselessChart(builder *kube.Builder) *kube.Chart {
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

	_, redisServiceName := chart.NewRedisStatefulsetWithOpts(namespace, kube.RedisOpts{
		AppPort: uint(redisPort),
	})

	env := []kube.EnvEntry{
		{Name: "REDIS_HOST", Value: kube.EnvValue{Value: fmt.Sprintf("%s.%s", redisServiceName, namespace)}},
		{Name: "REDIS_PORT", Value: kube.EnvValue{Value: fmt.Sprintf("%d", redisPort)}},
	}

	chart.NewDeployment(&kube.Deployment{
		Name:            name,
		Labels:          labels,
		Env:             env,
		Image:           appImage,
		ImagePullPolicy: "Always",
	})

	chart.NewIngress(&kube.Ingress{
		Name:   name,
		Labels: labels,
		Ingresses: []string{
			appIngress,
		},
		Port: uint(appPort),
	})

	chart.NewServiceMonitor(&kube.ServiceMonitor{
		Labels: labels,
	})

	return chart
}

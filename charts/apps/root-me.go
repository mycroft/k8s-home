package apps

import (
	kube "git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewRootMeChart(builder *kube.Builder) *kube.Chart {
	name := "root-me"

	namespace := name
	appName := name
	appImage := builder.RegisterContainerImage("registry.mkz.me/mycroft/scoreboard-api")
	appPort := 8080
	appIngress := "root-me.iop.cx"

	chart := builder.NewChart(name)
	chart.NewNamespace(namespace)

	labels := map[string]string{
		"app.kubernetes.io/name":      appName,
		"app.kubernetes.io/component": "api",
	}

	chart.NewDeployment(&kube.Deployment{
		Name:            name,
		Labels:          labels,
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

	return chart
}

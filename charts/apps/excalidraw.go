package apps

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func NewExcalidrawChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	name := "excalidraw"

	chart := builder.NewChart(name)
	chart.NewNamespace(name)

	labels := map[string]string{
		"app.kubernetes.io/name":      name,
		"app.kubernetes.io/component": "main",
	}

	chart.NewDeployment(&kubehelpers.Deployment{
		Name:            name,
		Labels:          labels,
		Image:           builder.RegisterContainerImage("excalidraw/excalidraw"),
		ImagePullPolicy: "Always",
	})

	chart.NewIngress(&kubehelpers.Ingress{
		Name:   name,
		Labels: labels,
		Port:   80,
		Ingresses: []string{
			"excalidraw.services.mkz.me",
		},
	})

	return chart
}

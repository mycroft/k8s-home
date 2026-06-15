package observability

import "git.mkz.me/mycroft/k8s-home/internal/kubehelpers"

func NewCustomMonitoring(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "monitoring"

	chart := builder.NewChart("custom-monitoring")

	// Custom monitors
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "garage-monitoring")

	// Scrape the Garage backup target's metrics endpoint, which is not reachable via a
	// Kubernetes Service since it lives outside the cluster.
	kubehelpers.CreateScrapeConfig(chart, kubehelpers.ScrapeTarget{
		Name:      "garage-monitoring",
		Namespace: namespace,
		Target:    "moonstone.lan.mkz.me:3903",
		BearerToken: &kubehelpers.ScrapeTargetBearerToken{
			SecretName: "garage-monitoring",
			SecretKey:  "metrics_token",
		},
	})

	// Scrape moonstone's node-exporter directly, since it isn't backed by a Kubernetes Service.
	kubehelpers.CreateScrapeConfig(chart, kubehelpers.ScrapeTarget{
		Name:      "moonstone-node-exporter",
		Namespace: namespace,
		Target:    "moonstone.lan.mkz.me:9100",
	})

	kubehelpers.CreateScrapeConfig(chart, kubehelpers.ScrapeTarget{
		Name:      "moonstone-zml-smi-node-exporter",
		Namespace: namespace,
		Target:    "moonstone.lan.mkz.me:9101",
	})

	// Scrape glitter's node-exporter directly, since it isn't backed by a Kubernetes Service.
	kubehelpers.CreateScrapeConfig(chart, kubehelpers.ScrapeTarget{
		Name:      "glitter-node-exporter",
		Namespace: namespace,
		Target:    "glitter.lan.mkz.me:9100",
	})

	return chart
}

package apps

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

// NewAgentNexusChart deploys agent-nexus, the LLM agent service that reviews
// pull requests (https://git.mkz.me/mycroft/agent-nexus), from its Helm
// chart. Values are in configs/agent-nexus.yaml.
//
// Its credentials live in Vault at secret/namespaces/agent-nexus/agent-nexus:
// api-keys, gitea-token, llm-api-key, and github-app-private-key once GitHub
// is enabled.
func NewAgentNexusChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "agent-nexus"
	repositoryName := "homelab-oci" // declared by NewHomelabOCIHelmRepositoryChart
	chartName := "agent-nexus"
	releaseName := "agent-nexus"
	appIngress := "agent-nexus.services.mkz.me"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "agent-nexus")

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	// The event API authenticates callers with API keys, so the ingress
	// has no auth middleware. It only routes the service's http port;
	// metrics stay on the separate metrics port.
	chart.NewIngress(&kubehelpers.Ingress{
		Name:        "agent-nexus",
		Ingresses:   []string{appIngress},
		ServiceName: releaseName, // created by the helm chart
	})

	return chart
}

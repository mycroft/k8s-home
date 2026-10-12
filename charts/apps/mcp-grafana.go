package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

// NewMCPGrafanaChart deploys mcp-grafana (https://github.com/grafana/mcp-grafana),
// through which agent-nexus reads Prometheus, Loki, Tempo and the dashboards
// via kube-prometheus-stack's Grafana. Values are in configs/mcp-grafana.yaml.
//
// Its credentials live in Vault at secret/namespaces/mcp-grafana/mcp-grafana:
// grafana-token, a token of the Grafana service account mcp-grafana (Viewer
// role, created by hand), and server-token, which callers must send as a
// bearer token. agent-nexus holds a copy of the latter as mcp-grafana-token.
//
// The server runs with --disable-write, has no ingress, and the
// NetworkPolicy only lets agent-nexus in.
func NewMCPGrafanaChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "mcp-grafana"
	repositoryName := "grafana-community"
	chartName := "grafana-mcp"
	releaseName := "mcp-grafana"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "mcp-grafana")

	chart.CreateHelmRepository(
		repositoryName,
		"https://grafana-community.github.io/helm-charts",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	newMCPGrafanaNetworkPolicy(chart, namespace, chartName, releaseName)

	return chart
}

// newMCPGrafanaNetworkPolicy admits agent-nexus to the MCP endpoint.
func newMCPGrafanaNetworkPolicy(chart *kubehelpers.Chart, namespace, chartName, releaseName string) {
	port := 8000.

	k8s.NewKubeNetworkPolicy(
		chart.Cdk8sChart,
		jsii.String("mcp-grafana"),
		&k8s.KubeNetworkPolicyProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("mcp-grafana"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.NetworkPolicySpec{
				PodSelector: &k8s.LabelSelector{
					MatchLabels: &map[string]*string{
						"app.kubernetes.io/name":     jsii.String(chartName),
						"app.kubernetes.io/instance": jsii.String(releaseName),
					},
				},
				PolicyTypes: &[]*string{
					jsii.String("Ingress"),
				},
				Ingress: &[]*k8s.NetworkPolicyIngressRule{
					{
						From: &[]*k8s.NetworkPolicyPeer{
							{
								NamespaceSelector: &k8s.LabelSelector{
									MatchLabels: &map[string]*string{
										"kubernetes.io/metadata.name": jsii.String("agent-nexus"),
									},
								},
								PodSelector: &k8s.LabelSelector{
									MatchLabels: &map[string]*string{
										"app.kubernetes.io/name": jsii.String("agent-nexus"),
									},
								},
							},
						},
						Ports: &[]*k8s.NetworkPolicyPort{
							{Port: k8s.IntOrString_FromNumber(&port), Protocol: jsii.String("TCP")},
						},
					},
				},
			},
		},
	)
}

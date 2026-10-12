package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

// NewKubernetesMCPServerChart deploys kubernetes-mcp-server
// (https://github.com/containers/kubernetes-mcp-server), through which
// agent-nexus reads the live state of the cluster. Its ServiceAccount, created
// by the chart, reads everything but Secrets. Values, including its read-only
// configuration, are in configs/kubernetes-mcp-server.yaml.tmpl.
//
// Upstream publishes every release as chart 0.1.0, re-pushed over the
// previous one: updates come through the image version in versions.yaml.
//
// The server has no authentication of its own and reads the cluster with its
// ServiceAccount: it has no ingress, and the NetworkPolicy only lets
// agent-nexus and Prometheus reach it.
func NewKubernetesMCPServerChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "kubernetes-mcp-server"
	repositoryName := "kubernetes-mcp-server"
	chartName := "kubernetes-mcp-server"
	releaseName := "kubernetes-mcp-server"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repositoryName,
		"oci://ghcr.io/containers/charts",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	chart.NewReadAllButSecretsClusterRole(releaseName) // the chart names its ServiceAccount after the release

	newKubernetesMCPServerNetworkPolicy(chart, namespace, releaseName)

	return chart
}

// newKubernetesMCPServerNetworkPolicy admits agent-nexus to the MCP endpoint
// and Prometheus to /metrics, both on the server's single port.
func newKubernetesMCPServerNetworkPolicy(chart *kubehelpers.Chart, namespace, releaseName string) {
	port := 8080.
	from := func(namespace, name string) *k8s.NetworkPolicyPeer {
		return &k8s.NetworkPolicyPeer{
			NamespaceSelector: &k8s.LabelSelector{
				MatchLabels: &map[string]*string{"kubernetes.io/metadata.name": jsii.String(namespace)},
			},
			PodSelector: &k8s.LabelSelector{
				MatchLabels: &map[string]*string{"app.kubernetes.io/name": jsii.String(name)},
			},
		}
	}

	k8s.NewKubeNetworkPolicy(
		chart.Cdk8sChart,
		jsii.String("kubernetes-mcp-server"),
		&k8s.KubeNetworkPolicyProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("kubernetes-mcp-server"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.NetworkPolicySpec{
				PodSelector: &k8s.LabelSelector{
					MatchLabels: &map[string]*string{
						"app.kubernetes.io/name":     jsii.String("kubernetes-mcp-server"),
						"app.kubernetes.io/instance": jsii.String(releaseName),
					},
				},
				PolicyTypes: &[]*string{
					jsii.String("Ingress"),
				},
				Ingress: &[]*k8s.NetworkPolicyIngressRule{
					{
						From: &[]*k8s.NetworkPolicyPeer{
							from("agent-nexus", "agent-nexus"),
							from("monitoring", "prometheus"),
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

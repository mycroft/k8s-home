package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

// NewAgentNexusChart deploys agent-nexus, the LLM agent service that reviews
// pull requests (https://git.mkz.me/mycroft/agent-nexus), from its Helm
// chart. Values are in configs/agent-nexus.yaml.
//
// Its credentials live in Vault at secret/namespaces/agent-nexus/agent-nexus:
// api-keys, gitea-token, llm-api-key, and github-app-private-key once GitHub
// is enabled.
//
// The jobs UI has no authentication of its own: it is served on its own port
// and host, behind an authentik proxy outpost. The authentik side is manual
// (see docs/authentik.md): a forward auth proxy provider for
// https://agent-nexus-ui.services.mkz.me, its application, and an outpost
// named "outpost-for-agent-nexus". authentik then creates the
// ak-outpost-outpost-for-agent-nexus middleware referenced by the UI ingress;
// until it exists, Traefik refuses that route.
func NewAgentNexusChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "agent-nexus"
	repositoryName := "homelab-oci" // declared by NewHomelabOCIHelmRepositoryChart
	chartName := "agent-nexus"
	releaseName := "agent-nexus"
	appIngress := "agent-nexus.services.mkz.me"
	uiIngress := "agent-nexus-ui.services.mkz.me"

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

	chart.NewIngress(&kubehelpers.Ingress{
		Name:            "agent-nexus-ui",
		Ingresses:       []string{uiIngress},
		ServiceName:     releaseName,
		PortName:        "ui",
		CertificateName: "agent-nexus-ui-tls",
		Annotations: map[string]string{
			"traefik.ingress.kubernetes.io/router.middlewares": "authentik-ak-outpost-outpost-for-agent-nexus@kubernetescrd",
		},
	})

	newAgentNexusNetworkPolicy(chart, namespace, releaseName)

	return chart
}

// newAgentNexusNetworkPolicy keeps the unauthenticated UI port reachable only
// through Traefik, where the authentik middleware applies, and lets
// Prometheus scrape the metrics port.
func newAgentNexusNetworkPolicy(chart *kubehelpers.Chart, namespace, releaseName string) {
	tcp := func(port float64) *k8s.NetworkPolicyPort {
		return &k8s.NetworkPolicyPort{Port: k8s.IntOrString_FromNumber(&port), Protocol: jsii.String("TCP")}
	}
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
		jsii.String("agent-nexus"),
		&k8s.KubeNetworkPolicyProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("agent-nexus"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.NetworkPolicySpec{
				PodSelector: &k8s.LabelSelector{
					MatchLabels: &map[string]*string{
						"app.kubernetes.io/name":     jsii.String("agent-nexus"),
						"app.kubernetes.io/instance": jsii.String(releaseName),
					},
				},
				PolicyTypes: &[]*string{
					jsii.String("Ingress"),
				},
				Ingress: &[]*k8s.NetworkPolicyIngressRule{
					{
						// The API (8080) and the UI (8081).
						From:  &[]*k8s.NetworkPolicyPeer{from("kube-system", "traefik")},
						Ports: &[]*k8s.NetworkPolicyPort{tcp(8080), tcp(8081)},
					},
					{
						From:  &[]*k8s.NetworkPolicyPeer{from("monitoring", "prometheus")},
						Ports: &[]*k8s.NetworkPolicyPort{tcp(9090)},
					},
				},
			},
		},
	)
}

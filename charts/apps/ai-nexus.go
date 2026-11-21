package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

// NewAINexusChart deploys ai-nexus, the AI factory control center that
// reviews pull requests with skill-based LLM passes
// (https://git.mkz.me/mycroft/ai-nexus), from its Helm chart. Values are in
// configs/ai-nexus.yaml.
//
// Its credentials live in Vault at secret/namespaces/ai-nexus/ai-nexus. Every
// key becomes an environment variable that the configuration references as
// ${KEY}: GITEA_TOKEN (the ai-nexus bot on git.mkz.me), GITEA_WEBHOOK_SECRET,
// GITHUB_APP_PRIVATE_KEY (the PEM of the GitHub App ai-nexus posts as on
// github.com), GITHUB_WEBHOOK_SECRET and LLAMA_API_KEY.
//
// Its runners (ADR-0016 in ai-nexus) run the review passes of repositories
// with engine = "claude-code" in Claude Code. Their one credential lives
// apart, in Vault at secret/namespaces/ai-nexus/ai-nexus-runner:
// CLAUDE_CODE_OAUTH_TOKEN, made with `claude setup-token`. The server never
// sees it, and the runners see nothing else.
//
// The tools service (ADR-0017) serves the web tools to the runners, for
// passes holding a grant the server signs with AI_NEXUS_GRANT_KEY (an
// Ed25519 private key PEM, in the server's Vault entry above); it checks
// grants with the public key in configs/ai-nexus.yaml. It reads the forge
// with a token of its own that can only read, in Vault at
// secret/namespaces/ai-nexus/ai-nexus-tools: GITEA_TOKEN.
//
// ai-nexus serves webhooks and its dashboard on the same port, so the split
// is by host and path: the webhook ingress routes /webhooks only, without
// auth (deliveries are HMAC-signed); the dashboard ingress routes its own host
// behind an authentik proxy outpost. The authentik side is manual (see
// docs/authentik.md): a forward auth proxy provider for
// https://ai-nexus-ui.services.mkz.me, its application, and an outpost named
// "outpost-for-ai-nexus". authentik then creates the
// ak-outpost-outpost-for-ai-nexus middleware referenced by the dashboard
// ingress; until it exists, Traefik refuses that route.
func NewAINexusChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "ai-nexus"
	repositoryName := "homelab-oci" // declared by NewHomelabOCIHelmRepositoryChart
	chartName := "ai-nexus"
	releaseName := "ai-nexus"
	webhookIngress := "ai-nexus.services.mkz.me"
	dashboardIngress := "ai-nexus-ui.services.mkz.me"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "ai-nexus")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "ai-nexus-runner")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "ai-nexus-tools")

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	// Forges deliver webhooks here; nothing else of the service is routed
	// on this host.
	chart.NewIngress(&kubehelpers.Ingress{
		Name:        "ai-nexus-webhooks",
		Ingresses:   []string{webhookIngress},
		ServiceName: releaseName, // created by the helm chart
		Path:        "/webhooks",
	})

	chart.NewIngress(&kubehelpers.Ingress{
		Name:            "ai-nexus-ui",
		Ingresses:       []string{dashboardIngress},
		ServiceName:     releaseName,
		CertificateName: "ai-nexus-ui-tls",
		Annotations: map[string]string{
			"traefik.ingress.kubernetes.io/router.middlewares": "authentik-ak-outpost-outpost-for-ai-nexus@kubernetescrd",
		},
	})

	newAINexusNetworkPolicy(chart, namespace, releaseName)
	newAINexusRunnerNetworkPolicy(chart, namespace, releaseName)
	newAINexusToolsNetworkPolicy(chart, namespace, releaseName)

	return chart
}

// aiNexusRunnerPort is the runner API's port on the server pod, set by the
// chart's runners.port; aiNexusToolsPort is the tools service's, set by
// tools.port.
const (
	aiNexusRunnerPort = 8081
	aiNexusToolsPort  = 8082
)

// newAINexusNetworkPolicy admits only Traefik, where the ingresses' path and
// authentik middleware apply, and Prometheus, which scrapes /metrics from the
// same port; and the runners, on the runner API's port only.
func newAINexusNetworkPolicy(chart *kubehelpers.Chart, namespace, releaseName string) {
	port := float64(8080)
	runnerPort := float64(aiNexusRunnerPort)
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
		jsii.String("ai-nexus"),
		&k8s.KubeNetworkPolicyProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("ai-nexus"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.NetworkPolicySpec{
				PodSelector: &k8s.LabelSelector{
					MatchLabels: &map[string]*string{
						"app.kubernetes.io/name":     jsii.String("ai-nexus"),
						"app.kubernetes.io/instance": jsii.String(releaseName),
					},
				},
				PolicyTypes: &[]*string{
					jsii.String("Ingress"),
				},
				Ingress: &[]*k8s.NetworkPolicyIngressRule{
					{
						From: &[]*k8s.NetworkPolicyPeer{
							from("kube-system", "traefik"),
							from("monitoring", "prometheus"),
						},
						Ports: &[]*k8s.NetworkPolicyPort{
							{Port: k8s.IntOrString_FromNumber(&port), Protocol: jsii.String("TCP")},
						},
					},
					{
						From: &[]*k8s.NetworkPolicyPeer{
							{PodSelector: aiNexusRunnerSelector(releaseName)},
						},
						Ports: &[]*k8s.NetworkPolicyPort{
							{Port: k8s.IntOrString_FromNumber(&runnerPort), Protocol: jsii.String("TCP")},
						},
					},
				},
			},
		},
	)
}

func aiNexusRunnerSelector(releaseName string) *k8s.LabelSelector {
	return &k8s.LabelSelector{
		MatchLabels: &map[string]*string{
			"app.kubernetes.io/name":     jsii.String("ai-nexus-runner"),
			"app.kubernetes.io/instance": jsii.String(releaseName),
		},
	}
}

// newAINexusRunnerNetworkPolicy closes the runners to every connection in,
// and lets them out only to the runner API, the tools service, the
// cluster's DNS and public addresses (Anthropic's API, and what a review
// session fetches): not the cluster's other services, the LAN or the
// tailnet.
func newAINexusRunnerNetworkPolicy(chart *kubehelpers.Chart, namespace, releaseName string) {
	runnerPort := float64(aiNexusRunnerPort)
	toolsPort := float64(aiNexusToolsPort)
	dnsPort := float64(53)
	httpPort := float64(80)
	httpsPort := float64(443)
	tcp := func(port *float64) *k8s.NetworkPolicyPort {
		return &k8s.NetworkPolicyPort{Port: k8s.IntOrString_FromNumber(port), Protocol: jsii.String("TCP")}
	}

	k8s.NewKubeNetworkPolicy(
		chart.Cdk8sChart,
		jsii.String("ai-nexus-runner"),
		&k8s.KubeNetworkPolicyProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("ai-nexus-runner"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.NetworkPolicySpec{
				PodSelector: aiNexusRunnerSelector(releaseName),
				PolicyTypes: &[]*string{
					jsii.String("Ingress"),
					jsii.String("Egress"),
				},
				Egress: &[]*k8s.NetworkPolicyEgressRule{
					{
						To: &[]*k8s.NetworkPolicyPeer{
							{PodSelector: &k8s.LabelSelector{
								MatchLabels: &map[string]*string{
									"app.kubernetes.io/name":     jsii.String("ai-nexus"),
									"app.kubernetes.io/instance": jsii.String(releaseName),
								},
							}},
						},
						Ports: &[]*k8s.NetworkPolicyPort{tcp(&runnerPort)},
					},
					{
						To:    &[]*k8s.NetworkPolicyPeer{{PodSelector: aiNexusToolsSelector(releaseName)}},
						Ports: &[]*k8s.NetworkPolicyPort{tcp(&toolsPort)},
					},
					{
						To: &[]*k8s.NetworkPolicyPeer{
							{
								NamespaceSelector: &k8s.LabelSelector{
									MatchLabels: &map[string]*string{"kubernetes.io/metadata.name": jsii.String("kube-system")},
								},
								PodSelector: &k8s.LabelSelector{
									MatchLabels: &map[string]*string{"k8s-app": jsii.String("kube-dns")},
								},
							},
						},
						Ports: &[]*k8s.NetworkPolicyPort{
							{Port: k8s.IntOrString_FromNumber(&dnsPort), Protocol: jsii.String("UDP")},
							tcp(&dnsPort),
						},
					},
					{
						To: &[]*k8s.NetworkPolicyPeer{
							{IpBlock: &k8s.IpBlock{
								Cidr: jsii.String("0.0.0.0/0"),
								Except: &[]*string{
									jsii.String("10.0.0.0/8"),
									jsii.String("172.16.0.0/12"),
									jsii.String("192.168.0.0/16"),
									jsii.String("100.64.0.0/10"),
									jsii.String("169.254.0.0/16"),
								},
							}},
						},
						Ports: &[]*k8s.NetworkPolicyPort{tcp(&httpPort), tcp(&httpsPort)},
					},
				},
			},
		},
	)
}

func aiNexusToolsSelector(releaseName string) *k8s.LabelSelector {
	return &k8s.LabelSelector{
		MatchLabels: &map[string]*string{
			"app.kubernetes.io/name":     jsii.String("ai-nexus-tools"),
			"app.kubernetes.io/instance": jsii.String(releaseName),
		},
	}
}

// newAINexusToolsNetworkPolicy admits only the runners, on the tools port,
// and lets the tools service out only to the cluster's DNS and public
// addresses: the forge (git.mkz.me resolves to one) and the chart
// repositories, registries and pages the web tools read.
func newAINexusToolsNetworkPolicy(chart *kubehelpers.Chart, namespace, releaseName string) {
	toolsPort := float64(aiNexusToolsPort)
	dnsPort := float64(53)
	httpPort := float64(80)
	httpsPort := float64(443)
	tcp := func(port *float64) *k8s.NetworkPolicyPort {
		return &k8s.NetworkPolicyPort{Port: k8s.IntOrString_FromNumber(port), Protocol: jsii.String("TCP")}
	}

	k8s.NewKubeNetworkPolicy(
		chart.Cdk8sChart,
		jsii.String("ai-nexus-tools"),
		&k8s.KubeNetworkPolicyProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("ai-nexus-tools"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.NetworkPolicySpec{
				PodSelector: aiNexusToolsSelector(releaseName),
				PolicyTypes: &[]*string{
					jsii.String("Ingress"),
					jsii.String("Egress"),
				},
				Ingress: &[]*k8s.NetworkPolicyIngressRule{
					{
						From:  &[]*k8s.NetworkPolicyPeer{{PodSelector: aiNexusRunnerSelector(releaseName)}},
						Ports: &[]*k8s.NetworkPolicyPort{tcp(&toolsPort)},
					},
				},
				Egress: &[]*k8s.NetworkPolicyEgressRule{
					{
						To: &[]*k8s.NetworkPolicyPeer{
							{
								NamespaceSelector: &k8s.LabelSelector{
									MatchLabels: &map[string]*string{"kubernetes.io/metadata.name": jsii.String("kube-system")},
								},
								PodSelector: &k8s.LabelSelector{
									MatchLabels: &map[string]*string{"k8s-app": jsii.String("kube-dns")},
								},
							},
						},
						Ports: &[]*k8s.NetworkPolicyPort{
							{Port: k8s.IntOrString_FromNumber(&dnsPort), Protocol: jsii.String("UDP")},
							tcp(&dnsPort),
						},
					},
					{
						To: &[]*k8s.NetworkPolicyPeer{
							{IpBlock: &k8s.IpBlock{
								Cidr: jsii.String("0.0.0.0/0"),
								Except: &[]*string{
									jsii.String("10.0.0.0/8"),
									jsii.String("172.16.0.0/12"),
									jsii.String("192.168.0.0/16"),
									jsii.String("100.64.0.0/10"),
									jsii.String("169.254.0.0/16"),
								},
							}},
						},
						Ports: &[]*k8s.NetworkPolicyPort{tcp(&httpPort), tcp(&httpsPort)},
					},
				},
			},
		},
	)
}

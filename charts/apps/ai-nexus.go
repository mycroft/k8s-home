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
// ${KEY}: GITEA_TOKEN (the ai-nexus bot on git.mkz.me), GITEA_WEBHOOK_SECRET
// and LLAMA_API_KEY.
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

	return chart
}

// newAINexusNetworkPolicy admits only Traefik, where the ingresses' path and
// authentik middleware apply, and Prometheus, which scrapes /metrics from the
// same port.
func newAINexusNetworkPolicy(chart *kubehelpers.Chart, namespace, releaseName string) {
	port := float64(8080)
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
				},
			},
		},
	)
}

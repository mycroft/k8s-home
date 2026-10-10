package infra

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

// NewFluxCDUIChart deploys fluxcd-ui (https://git.mkz.me/mycroft/fluxcd-ui)
// behind an authentik proxy outpost.
//
// The authentik side is manual (see docs/authentik.md): a forward auth proxy
// provider for https://flux.services.mkz.me, its application, and an outpost
// named "outpost-for-fluxcd-ui". authentik then creates the
// ak-outpost-outpost-for-fluxcd-ui middleware referenced by the ingress; until
// it exists, Traefik refuses the route.
func NewFluxCDUIChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "fluxcd-ui"
	repositoryName := "homelab-oci"
	chartName := "fluxcd-ui"
	releaseName := "fluxcd-ui"
	appIngress := "flux.services.mkz.me"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	// repositoryName is declared by NewHomelabOCIHelmRepositoryChart.
	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	// The UI trusts the user and groups headers set by the outpost: only
	// Traefik may reach the pod, so they can't be forged from inside the
	// cluster.
	podPort := 8080.
	k8s.NewKubeNetworkPolicy(
		chart.Cdk8sChart,
		jsii.String("fluxcd-ui"),
		&k8s.KubeNetworkPolicyProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("fluxcd-ui"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.NetworkPolicySpec{
				PodSelector: &k8s.LabelSelector{
					MatchLabels: &map[string]*string{
						"app.kubernetes.io/name":     jsii.String("fluxcd-ui"),
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
										"kubernetes.io/metadata.name": jsii.String("kube-system"),
									},
								},
								PodSelector: &k8s.LabelSelector{
									MatchLabels: &map[string]*string{
										"app.kubernetes.io/name": jsii.String("traefik"),
									},
								},
							},
						},
						Ports: &[]*k8s.NetworkPolicyPort{
							{
								Port:     k8s.IntOrString_FromNumber(&podPort),
								Protocol: jsii.String("TCP"),
							},
						},
					},
				},
			},
		},
	)

	chart.NewIngress(&kubehelpers.Ingress{
		Name:        "fluxcd-ui",
		Ingresses:   []string{appIngress},
		ServiceName: releaseName, // created by the helm chart
		Annotations: map[string]string{
			"traefik.ingress.kubernetes.io/router.middlewares": "authentik-ak-outpost-outpost-for-fluxcd-ui@kubernetescrd",
		},
	})

	return chart
}

package security

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

// NewTrivyOperatorDashboardChart deploys trivy-operator-dashboard
// (https://github.com/raoulx24/trivy-operator-dashboard), a read-only web UI
// over the trivy-operator reports, behind an authentik proxy outpost.
//
// The authentik side is manual (see docs/authentik.md): a forward auth proxy
// provider for https://trivy.services.mkz.me, its application, and an outpost
// named "outpost-for-trivy-operator-dashboard". authentik then creates the
// ak-outpost-outpost-for-trivy-operator-dashboard middleware referenced by the
// ingress; until it exists, Traefik refuses the route.
func NewTrivyOperatorDashboardChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "trivy-operator-dashboard"
	repositoryName := "raoulx24"
	chartName := "trivy-operator-dashboard"
	releaseName := "trivy-operator-dashboard"
	appIngress := "trivy.services.mkz.me"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repositoryName,
		"oci://ghcr.io/raoulx24/charts",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	// The dashboard has no authentication of its own and shows every finding
	// in the cluster: only Traefik, which applies the forward auth, may reach
	// the pod.
	podPort := 8900.
	k8s.NewKubeNetworkPolicy(
		chart.Cdk8sChart,
		jsii.String("trivy-operator-dashboard"),
		&k8s.KubeNetworkPolicyProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("trivy-operator-dashboard"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.NetworkPolicySpec{
				PodSelector: &k8s.LabelSelector{
					MatchLabels: &map[string]*string{
						"app.kubernetes.io/name":     jsii.String("trivy-operator-dashboard"),
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
		Name:        "trivy-operator-dashboard",
		Ingresses:   []string{appIngress},
		ServiceName: releaseName, // created by the helm chart
		Annotations: map[string]string{
			"traefik.ingress.kubernetes.io/router.middlewares": "authentik-ak-outpost-outpost-for-trivy-operator-dashboard@kubernetescrd",
		},
	})

	return chart
}

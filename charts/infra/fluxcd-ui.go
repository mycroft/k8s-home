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
	repositoryName := "mycroft-oci"
	chartName := "fluxcd-ui"
	releaseName := "fluxcd-ui"
	appIngress := "flux.services.mkz.me"

	// authentik group allowed to suspend, resume and reconcile. The UI
	// prefixes group names with "fluxcd-ui:" before checking RBAC.
	operatorGroup := "authentik Admins"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repositoryName,
		"oci://registry.mkz.me/mycroft/charts",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	// fluxcd-ui-operator is created by the helm chart.
	k8s.NewKubeClusterRoleBinding(
		chart.Cdk8sChart,
		jsii.String("operators"),
		&k8s.KubeClusterRoleBindingProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String("fluxcd-ui-operators"),
			},
			Subjects: &[]*k8s.Subject{
				{
					ApiGroup: jsii.String("rbac.authorization.k8s.io"),
					Kind:     jsii.String("Group"),
					Name:     jsii.String("fluxcd-ui:" + operatorGroup),
				},
			},
			RoleRef: &k8s.RoleRef{
				ApiGroup: jsii.String("rbac.authorization.k8s.io"),
				Kind:     jsii.String("ClusterRole"),
				Name:     jsii.String("fluxcd-ui-operator"),
			},
		},
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

package infra

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewCapacitorChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "flux-system"
	repositoryName := "onechart"

	chart := builder.NewChart("capacitor")

	kubehelpers.CreateHelmRepository(
		chart.Cdk8sChart,
		repositoryName,
		"https://chart.onechart.dev",
	)

	kubehelpers.CreateHelmRelease(
		chart.Cdk8sChart,
		namespace,
		repositoryName,
		"onechart",
		"capacitor",
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
				namespace,
				repositoryName,
				"capacitor.yaml",
			),
		},
		nil,
	)

	k8s.NewKubeServiceAccount(
		chart.Cdk8sChart,
		jsii.String("sa"),
		&k8s.KubeServiceAccountProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("capacitor"),
			},
		},
	)

	k8s.NewKubeClusterRole(
		chart.Cdk8sChart,
		jsii.String("clusterrole"),
		&k8s.KubeClusterRoleProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("capacitor"),
			},
			Rules: &[]*k8s.PolicyRule{
				{
					ApiGroups: jsii.Strings(
						"networking.k8s.io",
						"apps",
						"",
					),
					Resources: jsii.Strings(
						"pods",
						"pods/log",
						"ingresses",
						"deployments",
						"services",
						"secrets",
						"events",
						"configmaps",
					),
					Verbs: jsii.Strings(
						"get",
						"watch",
						"list",
					),
				},
				{
					ApiGroups: jsii.Strings(
						"source.toolkit.fluxcd.io",
						"kustomize.toolkit.fluxcd.io",
						"helm.toolkit.fluxcd.io",
					),
					Resources: jsii.Strings(
						"gitrepositories",
						"ocirepositories",
						"buckets",
						"helmrepositories",
						"helmcharts",
						"kustomizations",
						"helmreleases",
					),
					Verbs: jsii.Strings(
						"get",
						"watch",
						"list",
						"patch",
					),
				},
			},
		},
	)

	k8s.NewKubeClusterRoleBinding(
		chart.Cdk8sChart,
		jsii.String("clusterrolebinding"),
		&k8s.KubeClusterRoleBindingProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String("capacitor"),
			},
			Subjects: &[]*k8s.Subject{
				{
					Kind:      jsii.String("ServiceAccount"),
					Name:      jsii.String("capacitor"),
					Namespace: jsii.String(namespace),
				},
			},
			RoleRef: &k8s.RoleRef{
				Kind:     jsii.String("ClusterRole"),
				Name:     jsii.String("capacitor"),
				ApiGroup: jsii.String("rbac.authorization.k8s.io"),
			},
		},
	)

	labels := map[string]*string{
		"app.kubernetes.io/component": jsii.String("capacitor"),
	}

	annotations := map[string]string{
		"traefik.ingress.kubernetes.io/redirect-entry-point": "https",
		"traefik.ingress.kubernetes.io/redirect-permanent":   "true",
		"traefik.ingress.kubernetes.io/router.middlewares":   "traefik-forward-auth-traefik-forward-auth@kubernetescrd",
	}

	podPort := 9000.
	k8s.NewKubeNetworkPolicy(
		chart.Cdk8sChart,
		jsii.String("capacitor"),
		&k8s.KubeNetworkPolicyProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("capacitor"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.NetworkPolicySpec{
				PodSelector: &k8s.LabelSelector{
					MatchLabels: &map[string]*string{
						"app.kubernetes.io/instance": jsii.String("capacitor"),
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

	k8s.NewKubeNetworkPolicy(
		chart.Cdk8sChart,
		jsii.String("acme-capacitor"),
		&k8s.KubeNetworkPolicyProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("acme-capacitor"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.NetworkPolicySpec{
				PodSelector: &k8s.LabelSelector{
					MatchLabels: &map[string]*string{
						"acme.cert-manager.io/http01-solver": jsii.String("true"),
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
					},
				},
			},
		},
	)

	kubehelpers.NewAppIngress(
		builder.Context,
		chart.Cdk8sChart,
		labels,
		namespace,
		9000,
		"capacitor.services.mkz.me", // fqdn
		"capacitor",                 // service name, created by the helm chart
		annotations,
	)

	return chart
}

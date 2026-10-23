package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

// NewReviewsChart deploys reviews, the dashboard of the GitHub pull requests
// waiting for your review (https://git.mkz.me/mycroft/reviews), from its Helm
// chart. Values are in configs/reviews.yaml.
//
// Users sign in with GitHub, through an OAuth App whose callback URL is
// https://reviews.iop.cx/auth/callback; its client ID is in the values. Its
// client secret and the key encrypting the session cookies live in Vault at
// secret/namespaces/reviews/reviews: github-client-secret and session-key
// (openssl rand -base64 32). The pod reads them at startup, so restart it
// after changing them.
func NewReviewsChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "reviews"
	repositoryName := "homelab-oci" // declared by NewHomelabOCIHelmRepositoryChart
	chartName := "reviews"
	releaseName := "reviews"
	appIngress := "reviews.iop.cx"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "reviews")

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	// The app authenticates users itself, with GitHub: no auth middleware.
	chart.NewIngress(&kubehelpers.Ingress{
		Name:        "reviews",
		Ingresses:   []string{appIngress},
		ServiceName: releaseName, // created by the helm chart
	})

	newReviewsNetworkPolicy(chart, namespace, releaseName)

	return chart
}

// newReviewsNetworkPolicy lets only Traefik reach the pod. Egress stays open:
// the app calls the GitHub API.
func newReviewsNetworkPolicy(chart *kubehelpers.Chart, namespace, releaseName string) {
	podPort := 8080.
	k8s.NewKubeNetworkPolicy(
		chart.Cdk8sChart,
		jsii.String("reviews"),
		&k8s.KubeNetworkPolicyProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("reviews"),
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.NetworkPolicySpec{
				PodSelector: &k8s.LabelSelector{
					MatchLabels: &map[string]*string{
						"app.kubernetes.io/name":     jsii.String("reviews"),
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
}

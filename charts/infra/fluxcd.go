package infra

import (
	fluxcd_ocirepositories_sourcetoolkitfluxcdio "git.mkz.me/mycroft/k8s-home/imports/ocirepositories_sourcetoolkitfluxcdio"
	podmonitor "git.mkz.me/mycroft/k8s-home/imports/podmonitor_monitoringcoreoscom"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewFluxCDChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "flux-system"

	chart := builder.NewChart("fluxcd")

	// Replicate secret from vault to flux-system namespace
	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "registry")

	// Add the OCI repository to the flux-system namespace
	fluxcd_ocirepositories_sourcetoolkitfluxcdio.NewOciRepository(
		chart.Cdk8sChart,
		jsii.String("oci-repository-flux"),
		&fluxcd_ocirepositories_sourcetoolkitfluxcdio.OciRepositoryProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name: jsii.String("flux-system"),
			},
			Spec: &fluxcd_ocirepositories_sourcetoolkitfluxcdio.OciRepositorySpec{
				Url:      jsii.String("oci://registry.mkz.me/k8s-home/k8s-home/manifests"),
				Interval: jsii.String("1m0s"),
				Ref: &fluxcd_ocirepositories_sourcetoolkitfluxcdio.OciRepositorySpecRef{
					Tag: jsii.String("latest"),
				},
				SecretRef: &fluxcd_ocirepositories_sourcetoolkitfluxcdio.OciRepositorySpecSecretRef{
					Name: jsii.String("registry"),
				},
			},
		},
	)

	podmonitor.NewPodMonitor(
		chart.Cdk8sChart,
		jsii.String("podmonitor-flux"),
		&podmonitor.PodMonitorProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("flux-system"),
				Namespace: jsii.String(namespace),
				Labels: &map[string]*string{
					"app.kubernetes.io/part-of":   jsii.String("flux"),
					"app.kubernetes.io/component": jsii.String("monitoring"),
					"release":                     jsii.String("prometheus"),
				},
			},
			Spec: &podmonitor.PodMonitorSpec{
				NamespaceSelector: &podmonitor.PodMonitorSpecNamespaceSelector{
					MatchNames: &[]*string{
						jsii.String("flux-system"),
					},
				},
				Selector: &podmonitor.PodMonitorSpecSelector{
					MatchExpressions: &[]*podmonitor.PodMonitorSpecSelectorMatchExpressions{
						{
							Key:      jsii.String("app"),
							Operator: jsii.String("In"),
							Values: &[]*string{
								jsii.String("helm-controller"),
								jsii.String("source-controller"),
								jsii.String("kustomize-controller"),
								jsii.String("notification-controller"),
								jsii.String("image-automation-controller"),
								jsii.String("image-reflector-controller"),
							},
						},
					},
				},
				PodMetricsEndpoints: &[]*podmonitor.PodMonitorSpecPodMetricsEndpoints{
					{
						Port: jsii.String("http-prom"),
						Relabelings: &[]*podmonitor.PodMonitorSpecPodMetricsEndpointsRelabelings{
							{
								SourceLabels: &[]*string{
									jsii.String("__meta_kubernetes_pod_phase"),
								},
								Action: podmonitor.PodMonitorSpecPodMetricsEndpointsRelabelingsAction_KEEP,
								Regex:  jsii.String("Running"),
							},
						},
					},
				},
			},
		},
	)

	return chart
}

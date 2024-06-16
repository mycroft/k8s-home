package infra

import (
	podmonitor "git.mkz.me/mycroft/k8s-home/imports/podmonitor_monitoringcoreoscom"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewFluxCDChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "flux-system"

	chart := builder.NewChart("fluxcd")

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

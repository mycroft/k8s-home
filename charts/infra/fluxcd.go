package infra

import (
	"context"

	podmonitor "git.mkz.me/mycroft/k8s-home/imports/podmonitor_monitoringcoreoscom"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewFluxCDChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	namespace := "flux-system"

	chart := cdk8s.NewChart(
		scope,
		jsii.String("fluxcd"),
		&cdk8s.ChartProps{},
	)

	podmonitor.NewPodMonitor(
		chart,
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

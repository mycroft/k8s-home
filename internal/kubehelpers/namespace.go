package kubehelpers

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/jsii-runtime-go"
)

// NewNamespace creates a Kubernetes Namespace resource in the chart and records it
// as the chart's canonical namespace, used by NewDeployment, NewIngress, and other
// helpers that need to know which namespace they're operating in.
//
// Each chart is expected to own exactly one namespace: calling this twice on the
// same chart panics to catch accidental misconfiguration at synth time.
func (chart *Chart) NewNamespace(name string) k8s.KubeNamespace {
	if chart.Namespace != "" {
		panic("can not overwrite this chart's namespace")
	}

	chart.Namespace = name

	return k8s.NewKubeNamespace(
		chart.Cdk8sChart,
		jsii.String(fmt.Sprintf("ns-%s", name)),
		&k8s.KubeNamespaceProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String(name),
			},
		},
	)
}

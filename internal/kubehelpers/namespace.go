package kubehelpers

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/jsii-runtime-go"
)

// NewNamespace creates a Namespace in Chart and returns its k8s.KubeNamespace object
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

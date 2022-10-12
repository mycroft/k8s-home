package k8s_helpers

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func NewNamespace(chart constructs.Construct, name string) k8s.KubeNamespace {
	return k8s.NewKubeNamespace(
		chart,
		jsii.String(fmt.Sprintf("ns-%s", name)),
		&k8s.KubeNamespaceProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String(name),
			},
		},
	)
}

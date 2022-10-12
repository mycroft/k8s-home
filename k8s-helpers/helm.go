package k8s_helpers

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/sourcetoolkitfluxcdio"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// CreateHelmRepository creates a HelmRepository into the flux-system namespace
// similar to:
// - helm repo add jetstack https://charts.jetstack.io
func CreateHelmRepository(chart constructs.Construct, name, url string) sourcetoolkitfluxcdio.HelmRepository {
	return sourcetoolkitfluxcdio.NewHelmRepository(
		chart,
		jsii.String(fmt.Sprintf("helm-repo-%s", name)),
		&sourcetoolkitfluxcdio.HelmRepositoryProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String(name),
				Namespace: jsii.String("flux-system"),
			},
			Spec: &sourcetoolkitfluxcdio.HelmRepositorySpec{
				Url:      jsii.String(url),
				Interval: jsii.String("1m0s"),
			},
		},
	)
}

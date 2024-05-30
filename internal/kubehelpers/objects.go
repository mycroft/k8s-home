package kubehelpers

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/bitnamicom"
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
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

func NewSealedSecret(chart constructs.Construct, namespace, name string, secrets map[string]*string) bitnamicom.SealedSecret {
	return bitnamicom.NewSealedSecret(
		chart,
		jsii.String("ss"),
		&bitnamicom.SealedSecretProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String(name),
				Namespace: jsii.String(namespace),
			},
			Spec: &bitnamicom.SealedSecretSpec{
				EncryptedData: &secrets,
			},
		},
	)
}

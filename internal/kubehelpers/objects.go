package kubehelpers

import (
	"git.mkz.me/mycroft/k8s-home/imports/bitnamicom"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

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

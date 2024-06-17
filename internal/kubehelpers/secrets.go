package kubehelpers

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/externalsecrets_externalsecretsio"
	"git.mkz.me/mycroft/k8s-home/imports/secretstore_externalsecretsio"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func CreateSecretStore(chart constructs.Construct, namespace string) {
	secretstore_externalsecretsio.NewSecretStoreV1Beta1(
		chart,
		jsii.String("secret-store"),
		&secretstore_externalsecretsio.SecretStoreV1Beta1Props{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("secretstore-vault"),
			},
			Spec: &secretstore_externalsecretsio.SecretStoreV1Beta1Spec{
				Provider: &secretstore_externalsecretsio.SecretStoreV1Beta1SpecProvider{
					Vault: &secretstore_externalsecretsio.SecretStoreV1Beta1SpecProviderVault{
						Server:  jsii.String("http://vault.vault:8200"),
						Path:    jsii.String("secret"),
						Version: secretstore_externalsecretsio.SecretStoreV1Beta1SpecProviderVaultVersion_V2,
						Auth: &secretstore_externalsecretsio.SecretStoreV1Beta1SpecProviderVaultAuth{
							Kubernetes: &secretstore_externalsecretsio.SecretStoreV1Beta1SpecProviderVaultAuthKubernetes{
								MountPath: jsii.String("kubernetes"),
								Role:      jsii.String("external-secrets"),
							},
						},
					},
				},
			},
		},
	)
}

func (chart *Chart) CreateSecretStore(namespace string) {
	CreateSecretStore(chart.Cdk8sChart, namespace)
}

func CreateExternalSecret(chart constructs.Construct, namespace, name string) {
	externalsecrets_externalsecretsio.NewExternalSecretV1Beta1(
		chart,
		jsii.String(fmt.Sprintf("es-%s", name)),
		&externalsecrets_externalsecretsio.ExternalSecretV1Beta1Props{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
			},
			Spec: &externalsecrets_externalsecretsio.ExternalSecretV1Beta1Spec{
				DataFrom: &[]*externalsecrets_externalsecretsio.ExternalSecretV1Beta1SpecDataFrom{
					{
						Extract: &externalsecrets_externalsecretsio.ExternalSecretV1Beta1SpecDataFromExtract{
							ConversionStrategy: jsii.String("Default"),
							Key:                jsii.String(fmt.Sprintf("secret/namespaces/%s/%s", namespace, name)),
						},
					},
				},
				RefreshInterval: jsii.String("15m"),
				SecretStoreRef: &externalsecrets_externalsecretsio.ExternalSecretV1Beta1SpecSecretStoreRef{
					Kind: jsii.String("SecretStore"),
					Name: jsii.String("secretstore-vault"),
				},
				Target: &externalsecrets_externalsecretsio.ExternalSecretV1Beta1SpecTarget{
					CreationPolicy: externalsecrets_externalsecretsio.ExternalSecretV1Beta1SpecTargetCreationPolicy_OWNER,
					DeletionPolicy: externalsecrets_externalsecretsio.ExternalSecretV1Beta1SpecTargetDeletionPolicy_DELETE,
					Immutable:      jsii.Bool(false),
					Name:           jsii.String(name),
					Template: &externalsecrets_externalsecretsio.ExternalSecretV1Beta1SpecTargetTemplate{
						EngineVersion: jsii.String("v2"),
						Type:          jsii.String("Opaque"),
					},
				},
			},
		},
	)
}

func (chart *Chart) CreateExternalSecret(namespace, name string) {
	CreateExternalSecret(chart.Cdk8sChart, namespace, name)
}

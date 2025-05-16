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
	secretstore_externalsecretsio.NewSecretStore(
		chart,
		jsii.String("secret-store"),
		&secretstore_externalsecretsio.SecretStoreProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("secretstore-vault"),
			},
			Spec: &secretstore_externalsecretsio.SecretStoreSpec{
				Provider: &secretstore_externalsecretsio.SecretStoreSpecProvider{
					Vault: &secretstore_externalsecretsio.SecretStoreSpecProviderVault{
						Server:  jsii.String("http://vault.vault:8200"),
						Path:    jsii.String("secret"),
						Version: secretstore_externalsecretsio.SecretStoreSpecProviderVaultVersion_V2,
						Auth: &secretstore_externalsecretsio.SecretStoreSpecProviderVaultAuth{
							Kubernetes: &secretstore_externalsecretsio.SecretStoreSpecProviderVaultAuthKubernetes{
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
	externalsecrets_externalsecretsio.NewExternalSecret(
		chart,
		jsii.String(fmt.Sprintf("es-%s", name)),
		&externalsecrets_externalsecretsio.ExternalSecretProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
			},
			Spec: &externalsecrets_externalsecretsio.ExternalSecretSpec{
				DataFrom: &[]*externalsecrets_externalsecretsio.ExternalSecretSpecDataFrom{
					{
						Extract: &externalsecrets_externalsecretsio.ExternalSecretSpecDataFromExtract{
							ConversionStrategy: externalsecrets_externalsecretsio.ExternalSecretSpecDataFromExtractConversionStrategy_DEFAULT,
							Key:                jsii.String(fmt.Sprintf("secret/namespaces/%s/%s", namespace, name)),
						},
					},
				},
				RefreshInterval: jsii.String("15m"),
				SecretStoreRef: &externalsecrets_externalsecretsio.ExternalSecretSpecSecretStoreRef{
					Kind: externalsecrets_externalsecretsio.ExternalSecretSpecSecretStoreRefKind_SECRET_STORE,
					Name: jsii.String("secretstore-vault"),
				},
				Target: &externalsecrets_externalsecretsio.ExternalSecretSpecTarget{
					CreationPolicy: externalsecrets_externalsecretsio.ExternalSecretSpecTargetCreationPolicy_OWNER,
					DeletionPolicy: externalsecrets_externalsecretsio.ExternalSecretSpecTargetDeletionPolicy_DELETE,
					Immutable:      jsii.Bool(false),
					Name:           jsii.String(name),
					Template: &externalsecrets_externalsecretsio.ExternalSecretSpecTargetTemplate{
						EngineVersion: externalsecrets_externalsecretsio.ExternalSecretSpecTargetTemplateEngineVersion_V2,
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

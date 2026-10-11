package kubehelpers

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/externalsecrets_externalsecretsio"
	password "git.mkz.me/mycroft/k8s-home/imports/password_generatorsexternalsecretsio"
	"git.mkz.me/mycroft/k8s-home/imports/secretstore_externalsecretsio"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

const (
	// Passwords are generated without symbols so they stay safe to embed in a DSN.
	generatedPasswordLength = 32
	generatedPasswordDigits = 6

	passwordGeneratorAPIVersion = "generators.external-secrets.io/v1alpha1"
)

// newPasswordGenerator declares a Password generator; an ExternalSecret
// referencing it by name materializes the password under the key password.
func newPasswordGenerator(chart constructs.Construct, id, namespace, name string) {
	password.NewPassword(
		chart,
		jsii.String(id),
		&password.PasswordProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String(name),
			},
			Spec: &password.PasswordSpec{
				Length:      jsii.Number(generatedPasswordLength),
				Digits:      jsii.Number(generatedPasswordDigits),
				Symbols:     jsii.Number(0),
				NoUpper:     jsii.Bool(false),
				AllowRepeat: jsii.Bool(true),
			},
		},
	)
}

// CreateGeneratedPasswordSecret creates a Secret named name holding a random
// password under the key password. It suits credentials that never leave the
// namespace, such as a broker's password shared by the broker and its
// client, so nothing has to be stored in Vault by hand.
//
// As for CNPG databases, refreshInterval 0 pins the password once generated,
// and the Secret is retained if the ExternalSecret is deleted: a service that
// only reads its password on first start would otherwise lock its clients out.
func CreateGeneratedPasswordSecret(chart constructs.Construct, namespace, name string) {
	newPasswordGenerator(chart, "password-"+name, namespace, name)

	externalsecrets_externalsecretsio.NewExternalSecret(
		chart,
		jsii.String("es-"+name),
		&externalsecrets_externalsecretsio.ExternalSecretProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String(name),
			},
			Spec: &externalsecrets_externalsecretsio.ExternalSecretSpec{
				RefreshInterval: jsii.String("0"),
				DataFrom: &[]*externalsecrets_externalsecretsio.ExternalSecretSpecDataFrom{
					{
						SourceRef: &externalsecrets_externalsecretsio.ExternalSecretSpecDataFromSourceRef{
							GeneratorRef: &externalsecrets_externalsecretsio.ExternalSecretSpecDataFromSourceRefGeneratorRef{
								ApiVersion: jsii.String(passwordGeneratorAPIVersion),
								Kind:       externalsecrets_externalsecretsio.ExternalSecretSpecDataFromSourceRefGeneratorRefKind_PASSWORD,
								Name:       jsii.String(name),
							},
						},
					},
				},
				Target: &externalsecrets_externalsecretsio.ExternalSecretSpecTarget{
					CreationPolicy: externalsecrets_externalsecretsio.ExternalSecretSpecTargetCreationPolicy_OWNER,
					DeletionPolicy: externalsecrets_externalsecretsio.ExternalSecretSpecTargetDeletionPolicy_RETAIN,
					Immutable:      jsii.Bool(false),
					Name:           jsii.String(name),
				},
			},
		},
	)
}

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
		jsii.String("es-"+name),
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

func CreateDockerConfigJSONSecret(chart constructs.Construct, namespace, name, registry string) {
	vaultKey := fmt.Sprintf("secret/namespaces/%s/%s", namespace, name)

	externalsecrets_externalsecretsio.NewExternalSecret(
		chart,
		jsii.String("es-"+name),
		&externalsecrets_externalsecretsio.ExternalSecretProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
			},
			Spec: &externalsecrets_externalsecretsio.ExternalSecretSpec{
				DataFrom: &[]*externalsecrets_externalsecretsio.ExternalSecretSpecDataFrom{
					{
						Extract: &externalsecrets_externalsecretsio.ExternalSecretSpecDataFromExtract{
							ConversionStrategy: externalsecrets_externalsecretsio.ExternalSecretSpecDataFromExtractConversionStrategy_DEFAULT,
							Key:                jsii.String(vaultKey),
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
						Type:          jsii.String("kubernetes.io/dockerconfigjson"),
						Data: &map[string]*string{
							".dockerconfigjson": jsii.String(fmt.Sprintf(`{
  "auths": {
    "%s": {
      "username": "{{ .username }}",
      "password": "{{ .password }}",
      "auth": "{{ printf "\x25s:\x25s" .username .password | b64enc }}"
    }
  }
}`, registry)),
						},
					},
				},
			},
		},
	)
}

func (chart *Chart) CreateExternalSecret(namespace, name string) {
	CreateExternalSecret(chart.Cdk8sChart, namespace, name)
}

func (chart *Chart) CreateDockerConfigJSONSecret(namespace, name, registry string) {
	CreateDockerConfigJSONSecret(chart.Cdk8sChart, namespace, name, registry)
}

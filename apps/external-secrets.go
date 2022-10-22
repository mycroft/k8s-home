package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/externalsecrets_externalsecretsio"
	"git.mkz.me/mycroft/k8s-home/imports/secretstore_externalsecretsio"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewExternalSecretsChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "external-secrets"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"external-secrets",
		"https://charts.external-secrets.io",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"external-secrets", // repo name
		"external-secrets", // chart name
		"external-secrets", // release name
		"0.6.0",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{},
		nil,
	)

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

	externalsecrets_externalsecretsio.NewExternalSecretV1Beta1(
		chart,
		jsii.String("es"),
		&externalsecrets_externalsecretsio.ExternalSecretV1Beta1Props{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("testaroo"),
			},
			Spec: &externalsecrets_externalsecretsio.ExternalSecretV1Beta1Spec{
				DataFrom: &[]*externalsecrets_externalsecretsio.ExternalSecretV1Beta1SpecDataFrom{
					{
						Extract: &externalsecrets_externalsecretsio.ExternalSecretV1Beta1SpecDataFromExtract{
							ConversionStrategy: jsii.String("Default"),
							Key:                jsii.String("secret/namespaces/external-secrets/testaroo"),
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
					Name:           jsii.String("testaroo"),
					Template: &externalsecrets_externalsecretsio.ExternalSecretV1Beta1SpecTargetTemplate{
						EngineVersion: jsii.String("v2"),
						Type:          jsii.String("Opaque"),
					},
				},
			},
		},
	)

	return chart
}

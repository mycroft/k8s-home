package apps

import (
	"log"
	"os"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewDexIdpChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "dex-idp"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"dex",
		"https://charts.dexidp.io",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"dex", // repo name
		"dex", // chart name
		"dex", // release name
		"0.12.1",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"dex-idp.yaml",
			),
		},
		nil,
	)

	// Create configuration
	// The configuration is stored in a secret, and secrets used are fetched from Vault using ExternalSecrets
	contents, err := os.ReadFile("configs/dex-config.yaml")
	if err != nil {
		log.Fatalf("Could not read config file for dex-idp: %v", err)
	}
	k8s.NewKubeSecret(
		chart,
		jsii.String("dex-config"),
		&k8s.KubeSecretProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("dex-config"), // referenced in helm chart config
			},
			Immutable: jsii.Bool(false),
			StringData: &map[string]*string{
				"config.yaml": jsii.String(string(contents)),
			},
		},
	)

	// Create ExternalSecrets
	k8s_helpers.CreateExternalSecret(chart, namespace, "static-admin")
	k8s_helpers.CreateExternalSecret(chart, namespace, "gitea")
	k8s_helpers.CreateExternalSecret(chart, namespace, "grafana-oidc-client")
	k8s_helpers.CreateExternalSecret(chart, namespace, "traefik-forward-auth-oidc")
	k8s_helpers.CreateExternalSecret(chart, namespace, "postgresql")

	return chart
}

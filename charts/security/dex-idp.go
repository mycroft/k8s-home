package security

import (
	"context"
	"log"
	"os"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewDexIdpChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	namespace := "dex-idp"
	repositoryName := "dex"
	chartName := "dex"
	releaseName := "dex"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)
	kubehelpers.CreateSecretStore(chart, namespace)

	kubehelpers.CreateHelmRepository(
		chart,
		"dex",
		"https://charts.dexidp.io",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName, // repo name
		chartName,      // chart name
		releaseName,    // release name
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
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
	kubehelpers.CreateExternalSecret(chart, namespace, "static-admin")
	kubehelpers.CreateExternalSecret(chart, namespace, "gitea")
	kubehelpers.CreateExternalSecret(chart, namespace, "grafana-oidc-client")
	kubehelpers.CreateExternalSecret(chart, namespace, "traefik-forward-auth-oidc")
	kubehelpers.CreateExternalSecret(chart, namespace, "postgresql")

	return chart
}

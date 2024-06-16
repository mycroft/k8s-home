package security

import (
	"log"
	"os"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewDexIdpChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "dex-idp"
	repositoryName := "dex"
	chartName := "dex"
	releaseName := "dex"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)

	chart.CreateHelmRepository(
		"dex",
		"https://charts.dexidp.io",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName, // repo name
		chartName,      // chart name
		releaseName,    // release name
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart.Cdk8sChart,
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
		chart.Cdk8sChart,
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
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "static-admin")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "gitea")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "grafana-oidc-client")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "traefik-forward-auth-oidc")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql")

	return chart
}

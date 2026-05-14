package security

import (
	"git.mkz.me/mycroft/k8s-home/imports/traefikio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewHarborScannerTrivy(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "harbor-scanner-trivy"
	namespace := appName

	repoName := "aqua"
	releaseName := "harbor-scanner-trivy"

	chart := builder.NewChart(appName)
	chart.NewNamespace(namespace)
	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "basic-auth-users")

	chart.CreateHelmRepository(
		repoName,
		"https://helm.aquasec.com",
	)

	values := map[string]*string{}
	labels := map[string]*string{}

	chart.NewRedisStatefulset(namespace)

	chart.CreateHelmRelease(
		namespace,
		repoName,
		appName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
		kubehelpers.WithHelmValues(values),
	)

	traefikio.NewMiddleware(
		chart.Cdk8sChart,
		jsii.String("basic-auth"),
		&traefikio.MiddlewareProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("basic-auth"),
				Namespace: jsii.String(namespace),
			},
			Spec: &traefikio.MiddlewareSpec{
				BasicAuth: &traefikio.MiddlewareSpecBasicAuth{
					Realm:  jsii.String("Trivy Scanner Authentication"),
					Secret: jsii.String("basic-auth-users"),
				},
			},
		},
	)

	kubehelpers.NewAppIngresses(
		builder.Context,
		chart.Cdk8sChart,
		labels,
		appName,
		uint(8080),
		[]string{"trivy-scanner.services.mkz.me"},
		"harbor-scanner-trivy",
		map[string]string{
			"traefik.ingress.kubernetes.io/router.entrypoints": "websecure",
			"traefik.ingress.kubernetes.io/router.middlewares": "harbor-scanner-trivy-basic-auth@kubernetescrd",
		},
		kubehelpers.AppIngressOption{
			PortName: "api-server",
		},
	)

	return chart
}

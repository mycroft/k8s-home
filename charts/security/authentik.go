package security

import (
	"context"
	"fmt"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewAuthentikChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	namespace := "authentik"
	repositoryName := "authentik"
	releaseName := "authentik"
	chartName := "authentik"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)
	kubehelpers.CreateSecretStore(chart, namespace)
	kubehelpers.CreateExternalSecret(chart, namespace, "postgresql")
	kubehelpers.CreateExternalSecret(chart, namespace, "authentik-secret")
	kubehelpers.CreateExternalSecret(chart, namespace, "mailrelay")

	_, redisServiceName := kubehelpers.NewRedisStatefulset(chart, namespace)
	_ = fmt.Sprintf("redis://%s:6379", redisServiceName)

	kubehelpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://charts.goauthentik.io",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName,
		chartName,
		releaseName,
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				releaseName,
				"authentik.yaml",
			),
		},
		nil,
	)

	return chart
}

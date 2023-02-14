package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewMinioOperator(scope constructs.Construct) cdk8s.Chart {
	namespace := "minio-operator"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"minio",
		"https://operator.min.io/",
	)

	k8s_helpers.CreateExternalSecret(chart, namespace, "root")

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,  // namespace
		"minio",    // repository name, same as above
		"operator", // the chart name
		"operator", // the release name
		"4.5.8",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"minio-operator.yaml",
			),
		},
		nil,
	)

	// XXX TODO: Create a service account secret
	/*
		apiVersion: v1
		kind: Secret
		metadata:
		  name: console-sa-secret
		  namespace: minio-operator
		  annotations:
		    kubernetes.io/service-account.name: console-sa
		type: kubernetes.io/service-account-token
	*/

	return chart
}

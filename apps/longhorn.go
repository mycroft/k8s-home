package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewLonghornChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "longhorn-system"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"longhorn",
		"https://charts.longhorn.io",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"longhorn", // repo name
		"longhorn", // chart name
		"longhorn", // release name
		"1.3.2",
		nil,
		nil,
		nil,
	)

	// Ingress requires basic-auth credentials
	k8s_helpers.CreateExternalSecret(chart, namespace, "basic-auth")

	k8s.NewKubeIngress(
		chart,
		jsii.String("ingress"),
		&k8s.KubeIngressProps{
			Metadata: &k8s.ObjectMeta{
				Annotations: &map[string]*string{
					"ingress.kubernetes.io/auth-type":    jsii.String("basic"),
					"ingress.kubernetes.io/auth-secret":  jsii.String("basic-auth"),
					"ingress.kubernetes.io/ssl-redirect": jsii.String("true"),
				},
				Namespace: jsii.String(namespace),
			},
			Spec: &k8s.IngressSpec{
				Rules: &[]*k8s.IngressRule{
					{
						Http: &k8s.HttpIngressRuleValue{
							Paths: &[]*k8s.HttpIngressPath{
								{
									PathType: jsii.String("Prefix"),
									Path:     jsii.String("/"),
									Backend: &k8s.IngressBackend{
										Service: &k8s.IngressServiceBackend{
											Name: jsii.String("longhorn-frontend"),
											Port: &k8s.ServiceBackendPort{
												Number: jsii.Number(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	)

	return chart
}

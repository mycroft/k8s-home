package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/imports/miniominio"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewMinio(scope constructs.Construct) cdk8s.Chart {
	namespace := "minio"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)

	k8s_helpers.CreateExternalSecret(chart, namespace, "minio-tenant")

	miniominio.NewTenantV2(
		chart,
		jsii.String("minio-storage"),
		&miniominio.TenantV2Props{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("minio-storage"),
				Namespace: jsii.String(namespace),
			},
			Spec: &miniominio.TenantV2Spec{
				CredsSecret: &miniominio.TenantV2SpecCredsSecret{
					Name: jsii.String("minio-tenant"),
				},
				Pools: &[]*miniominio.TenantV2SpecPools{
					{
						Servers:          jsii.Number(4),
						Name:             jsii.String("pool"),
						VolumesPerServer: jsii.Number(2),
						VolumeClaimTemplate: &miniominio.TenantV2SpecPoolsVolumeClaimTemplate{
							Metadata: &miniominio.TenantV2SpecPoolsVolumeClaimTemplateMetadata{
								Name: jsii.String("data"),
							},
							Spec: &miniominio.TenantV2SpecPoolsVolumeClaimTemplateSpec{
								StorageClassName: jsii.String("longhorn-crypto-global"),
								AccessModes: &[]*string{
									jsii.String("ReadWriteOnce"),
								},
								Resources: &miniominio.TenantV2SpecPoolsVolumeClaimTemplateSpecResources{
									Requests: &map[string]miniominio.TenantV2SpecPoolsVolumeClaimTemplateSpecResourcesRequests{
										"storage": miniominio.TenantV2SpecPoolsVolumeClaimTemplateSpecResourcesRequests_FromString(jsii.String("16Gi")),
									},
								},
							},
						},
						ContainerSecurityContext: &miniominio.TenantV2SpecPoolsContainerSecurityContext{
							RunAsUser:    jsii.Number(1000),
							RunAsGroup:   jsii.Number(1000),
							RunAsNonRoot: jsii.Bool(true),
						},
					},
				},
				RequestAutoCert: jsii.Bool(false),
				Env: &[]*miniominio.TenantV2SpecEnv{
					{
						Name:  jsii.String("MINIO_DOMAIN"),
						Value: jsii.String("minio-storage.services.mkz.me"),
					},
					{
						Name:  jsii.String("MINIO_SERVER_URL"),
						Value: jsii.String("https://minio-storage.services.mkz.me"),
					},
					{
						Name:  jsii.String("MINIO_BROWSER_REDIRECT_URL"),
						Value: jsii.String("https://minio-storage-console.services.mkz.me"),
					},
				},
			},
		},
	)

	// Create ingress for minio-storage & minio-storage-console
	// Using https://github.com/minio/operator/blob/master/examples/kustomization/tenant-letsencrypt/ingress.yaml

	annotations := map[string]*string{
		"kubernetes.io/ingress.class":                        jsii.String("traefik"),
		"cert-manager.io/cluster-issuer":                     jsii.String("letsencrypt-prod"),
		"traefik.ingress.kubernetes.io/redirect-entry-point": jsii.String("https"),
		"traefik.ingress.kubernetes.io/redirect-permanent":   jsii.String("true"),
	}

	storageIngress := "minio-storage.services.mkz.me"
	consoleIngress := "minio-storage-console.services.mkz.me"

	k8s.NewKubeIngress(
		chart,
		jsii.String("minio-storage-ingress"),
		&k8s.KubeIngressProps{
			Metadata: &k8s.ObjectMeta{
				Annotations: &annotations,
				Namespace:   jsii.String(namespace),
			},
			Spec: &k8s.IngressSpec{
				Rules: &[]*k8s.IngressRule{
					{
						Host: jsii.String(storageIngress),
						Http: &k8s.HttpIngressRuleValue{
							Paths: &[]*k8s.HttpIngressPath{
								{
									Backend: &k8s.IngressBackend{
										Service: &k8s.IngressServiceBackend{
											Name: jsii.String("minio"),
											Port: &k8s.ServiceBackendPort{
												Number: jsii.Number(80),
											},
										},
									},
									Path:     jsii.String("/"),
									PathType: jsii.String("Prefix"),
								},
							},
						},
					},
					{
						Host: jsii.String(consoleIngress),
						Http: &k8s.HttpIngressRuleValue{
							Paths: &[]*k8s.HttpIngressPath{
								{
									Backend: &k8s.IngressBackend{
										Service: &k8s.IngressServiceBackend{
											Name: jsii.String("minio-storage-console"),
											Port: &k8s.ServiceBackendPort{
												Number: jsii.Number(9090),
											},
										},
									},
									Path:     jsii.String("/"),
									PathType: jsii.String("Prefix"),
								},
							},
						},
					},
				},
				Tls: &[]*k8s.IngressTls{
					{
						Hosts: &[]*string{
							jsii.String(storageIngress),
							jsii.String(consoleIngress),
						},
						SecretName: jsii.String("secret-tls-www"),
					},
				},
			},
		},
	)

	return chart
}

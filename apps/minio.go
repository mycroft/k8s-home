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
	k8s_helpers.CreateExternalSecret(chart, namespace, "storage-configuration")

	miniominio.NewTenant(
		chart,
		jsii.String("minio-storage"),
		&miniominio.TenantProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("minio-storage"),
				Namespace: jsii.String(namespace),
			},
			Spec: &miniominio.TenantSpec{
				Configuration: &miniominio.TenantSpecConfiguration{
					Name: jsii.String("storage-configuration"),
				},
				Pools: &[]*miniominio.TenantSpecPools{
					{
						Servers:          jsii.Number(2),
						Name:             jsii.String("pool"),
						VolumesPerServer: jsii.Number(2),
						VolumeClaimTemplate: &miniominio.TenantSpecPoolsVolumeClaimTemplate{
							Metadata: &miniominio.TenantSpecPoolsVolumeClaimTemplateMetadata{
								Name: jsii.String("data"),
							},
							Spec: &miniominio.TenantSpecPoolsVolumeClaimTemplateSpec{
								StorageClassName: jsii.String("longhorn-crypto-global"),
								AccessModes: &[]*string{
									jsii.String("ReadWriteOnce"),
								},
								Resources: &miniominio.TenantSpecPoolsVolumeClaimTemplateSpecResources{
									Requests: &map[string]miniominio.TenantSpecPoolsVolumeClaimTemplateSpecResourcesRequests{
										"storage": miniominio.TenantSpecPoolsVolumeClaimTemplateSpecResourcesRequests_FromString(jsii.String("32Gi")),
									},
								},
							},
						},
						ContainerSecurityContext: &miniominio.TenantSpecPoolsContainerSecurityContext{
							RunAsUser:    jsii.Number(1000),
							RunAsGroup:   jsii.Number(1000),
							RunAsNonRoot: jsii.Bool(true),
						},
					},
				},
				Env: &[]*miniominio.TenantSpecEnv{
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
												Number: jsii.Number(443),
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
												Number: jsii.Number(9443),
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

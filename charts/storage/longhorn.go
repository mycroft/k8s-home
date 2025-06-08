package storage

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"

	"git.mkz.me/mycroft/k8s-home/imports/bitnamicom"
	"git.mkz.me/mycroft/k8s-home/imports/certificates_certmanagerio"
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/imports/longhornio"
	"git.mkz.me/mycroft/k8s-home/imports/servicemonitor_monitoringcoreoscom"
	"git.mkz.me/mycroft/k8s-home/imports/traefikio"

	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

const (
	cryptKeyCipherValue   = "AgAl+8NJCZvicrukusoTkrg+CYetiJb6485U0zX0v5cQcsnQq7F46T0GQpx6wnMMM4jUZqaYnm4GfBMtfH19Z8Pm3MXYbv/aIE3mzCKgvJcY0rd5aS7VgvLRuK7fBRI0z2iPEcVibWcF48AXrt7I18sYtZOrNVHIKmw22iIe2ILPdSIxP8yYDYdE5TVTPWiDGDO2xo041V83XkEoV82HWQmtUt4fx5iEyM+OrWi0v8tRqYO4yQxR3gpfVSLuGNajUri+2w/r+xpcaXso0pAmh4HX8AVGl8I7PSHhRw5R2gVo9DJHf0NFIpMu2YA6a0xsLTYpTFhBpQ7rWn2UJ4ZPAhsIzbad91R6WmRhv7AZwpYJmg+RUZc0YK+0HexhbTLdN2AH7bgee5OE31EPeYY5Z7e4uWw8OyrLaVv6+hQZl1jgykpSaDj6LsdacfEvyf8Opps7PSV+I966hfyDU+VkZkJbdyk4AYRW/cggMFt+fiyWvC2Vzbx2BavA2GYudmEcVqUVXz9bPwVddYqRL8wjRjaa+dxdrAOA4Vre1ILi0X3ylQFZ5bLB7u387+NhXTzg3AyM8WGopExF23GYMoKIZlzHHzHUPXnxHS2MpH7lBkuSN2QoCD7c5L4XzUly8pYP5rAYnH5BWLYti7YHa0S2XnX7EGpfrhxNb4uyriSb//mI4DLtgJUkQDYuW7X9YxpwbuTFPg5dgo5lmA5yD52a9O8="
	cryptKeyHashValue     = "AgCBnoQlm+emepq/LFzE18/7QNKtEH3mHTAGPbu9XEV4b/Pwy6pqunJfnTcq5zoKL5CV3PV4B1ZN0xuAxlf6Mc7xCLCCuvKIIh3YgKCy1aZNd4VymO2VVo1MBGdzbGfFKqjsCkhut4MYPVaoqk5r0wSLsZ8FAVz8mYiKyPANMJ+9HwLwk5laPMSrLfFlKRuEFsTanoSwZEQkLM50CME4b7+KpAFCa0Kn8AbMz7ue6ACXCvM7yRMWCLsAmPolrCT2EYXQKIasbI0m2+AIlCn27zwRFa5WIyGQ9ycdMm9WTpjvlGdEFxTU+8yp0EotGi8Q6H4k1IxGRnt7PG2y2wgn3Tp3boJ1f4WKN/Cdk2FxIisSaI8YBT2ARkQv63va4eXDUMivW58t3rPCARRxRH9tRZ26DTtIKZ5oIhGCIaU7m+XRo6scQNe9MSmYvaSmYueYO3p4sTcWOrUczAzEF2mqzTFhzTry7t/wpOcsHkoI4/Bwy0ZR2xmjpzN0DJmF0sb8IgvtAB3yQ+v1DhvZFfWi3oig4hyg3fn9UCCNWS1zYoAeuqTTwB32KqQYh8uzaThGZYt572KwFWtyWlqdKzOyeMA1I5Gs4rDNSvuYedQ9ff6fQwzgyCani1DNczErvFSXIf0GsARx5SP+rkDADxp6nIGoYCeE5m56QjRFTCoSRPCJZfPJUgT+IRyZQIvYWeKNYxxtZJn2mW8="
	cryptKeyProviderValue = "AgBriB/wtVOoyUzFyf84OYCNmKpBWE2JKxM/t5RoGknTqgXcboxSg/WH+uKwUQjl1oXTG85JdlMX5FXYEw5XY+UqCazc4efiXQTVvs1vSqYnm5r+cRVw2hJsX3pmZqp8zKbHmk5AWRL0kuuNi1E6x5zJ8mzleSA2ZfYGIRChstQ2c9vPES+D9IbDgyqOMtxGqGje7Ty7xj+zmKK3aREbJK66JWTpEYECvCPYS7/dZ9IzdXTNom2bJlybq9gpY0W2/Hz1OlgUU+VnkxKEQP1EsGFdkYyj8kOrdZTWX0doUbvkcoHH0lUhjBON5tIrdbXoQg5J8yaZ0bi18N6D+oEWCWYsVUg1hXwz/BZCNX8cLZBMg78HPEmdJ52Mn82fO7iiHCgm7LKoBPtQayJ/gehMVOwUS0DKgOIn5Ue8LRW5VwnKsKn3bTS0RtPHPny39039ycZ1jZteN97/b/OUwexFBT3fG5EhmS22iSBuf7kdWzWZuayjrWu93Blo/G+2M2O/Dzh5cJRaeWOwx0rC3o2WQUq8b9kdtxbebCtI1YTEqkg+jItvfgyedtNU+1Nu6S22C0Y6aAI2iERfPtVJXX04zM78HVEHlS5XETxOblPXwvOQplADSLUy3qeCXIVCYm3WZ0n+JrHPzey9tMPDPSf4NLZJi2yjt9u0UCCsjNGn/AgHtaCP+s63i1SakqcVbq6Lhy15Uf7mBsA="
	cryptKeySizeValue     = "AgAXmCIDpOmyPoVzKd4MMMsSZyqzxzId1gSOhEoMYV3E5DH1DP2HFv6s4MfrONQOFASjABliQbMa3P3VqH8FPqXCzyh+N6z5GNFRM9eeQHQ6TcffxFHRlA2nZ0QqoHsSvt0S0VbXIqj9lBsvIgJKEteQ8dRBWsF92BPHySk9GtlVMEAyGt2bNOyR/KO7wWQX1ses/rK5UfTrRd5OWfLj2KXfaD1dc5LmIK5xS0aBEPP8Oy7/G0SL7rVVe+1yViQqqN+I1WzXIsKG0pLEOOG3wxrHJBHe2JzJVpXp1plIWDEFJ1tEKQ4XExrF4iQ+mh88CbzIX6apARdpfvtqKgVtDMqnOKsZMWCQv4UP3/iYLwTTOCP02bFvQZXnfddhSv6k9kciV55wf0kTMKJcHq8Xzl3vX3tGQdHGUxmmBSoJJ+ztmuMc/uThHFZi9kovzA08V3gJ5EImDGju2snbRO2YDL+cWAGJ0jcJSdoTN4S4pjpbDEDrNqcEfkvSCT1SQl1Wv21NidSDw8PRDiGvwA+W2Hclwsmryj4Yv6M18VKkWCXWG5BfsflITI3/+tQt/t7C9W17fQewSL9O7IE9enppgaTLdr1vaASHe29E5ATljmFHGwse+jMxSlo/Ueo0ewRXn2ckBCoKtrdJs7f0ABLRQXjfvknU8ue+tUU+2KeoyTuxOIkefhXLNxG321CaW16/waVFvfY="
	cryptoKeyValueValue   = "AgC/qLeFS+oHqaDiiCJm3JyNvA5dk4MHOFxvGfo7YGpjB8YniVx54YzCP8EX47SAdX1D03DwYZVilSHe7W68US8bFSHRCUQE5PKC8ZcvkwJmwV53dwecdmLpwevUBY4KZ8mwT9CAs6g+10eQ2fOF9ydZB1QnOrwV2Tkp/tPoucswHSV0FbZmBWnEq96VZjD1GEVinPVXiGnmbAsoVWBBPqFNRb3TqvaGb6OE3W4SUzrZ4Vo/+rNag6NDd8xQx2gshl1M5xKwar33zFjmAKkKWPA8pDsKtenBSoxxGz/eENlZMvPJLksi2uaBc/U/3ModTGfdM4uuMTM46mLUvP1aiMh9CbwOXvzihnUKQ3E4cF1i8ZQTYkoTYND78l2iWjkurGCuayACTrE0CGJN3YqyyZIj3a8Qaa3KyghNkKbP5RAE8McnzcRBxbQi3OUXIQMqnOqmrT4uzpmsMahtce6ib4pYURunayOEajh+YmTzYWv2rrq3poEk/r+wP6n3SVZR6pj2PoQwHN1HMYIrQqDwwIpgvoBmj3OQ1mtlpdiIrIj2NYa9K0jRKa0L6WWb/RYdNI+z4p/YUJ3btn4ZA2+AjO90ecXGfkB+MKwO4QBdiljc25nugIfRuhmQ690EPojJ1DFL4HU+HSpjXLAtNpxvZE6Up5pUjniC9iDUHq1xMwCynA3RGg9DJik8MyHbDwtkW17jI328WTQm5QHvqb/ZhXJA2ukqjtjjCL5k4Y2DxYW9XDvPlteb3wKXz+rTAfzAaLnBWaqrBJ2BSZ8kcJdKEJAO"
	cryptoPBKDFValue      = "AgDKQtI+u+mVmrhK6vZi8PLd/g0hwu6n6132NRDN0pjf1TIjonIxvY7w/VDcYGdkz/a7FGQrs9sc8gNN59gJrEAuXNiYU1fxICdcrrvZIqrNvJKgOnWNbBxTe7vV2E7goYc238YHdJudIYHeas+9KB5AGHj5YyRfrjE8jMMG03PNykhIaG8ZE/x9lxs694gQNHSCmaquvRkDjtVbRQ5J08nAiUobLkNHprR+Qzz5lifPw6SpY3pUvljFxYmmpOL18t9F981bBC9T7B1Xt4+YL4oBCq8u7gWeWuPzEr1sCgrte9lBNmprEf+Wady+nV3wC+YQtAi3gjcWC1ygTlnrBU3rtosf6frouQbJ7Y7cIdKDpv1iQIFS2YTPJaUixCsNhxmhi6IrMrczwDOBunGA9ZIuj9oL2g0ObhRePixqNGtOvDvQeZ/BvoH35GjQNnSVhUiyKqdlQ0JgD7iFonTIa3IS9sTd26WS+ah31Aq66hXIC05QbPMpaOMEjx9nktMblKDF1hNC97d46vEShai5Qlbg1BAORZ2Otzb4wpFsKmP+ZlPtWiKSn1ypTnQmIgvnrM9/SkrhJkxt+SD6UpSe+19CbgzonC+CaU2fylOzk3Qr3bvY6Gc6ALaYvxI5QhlaNuVsrWW6IbTGIMH2ckkb3vBZ9kM03b1/UD/uMcHIw8K5G/KyLbDLX580Au8kmk7/IVyDUdPWpeKy"
)

func NewLonghornChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "longhorn-system"
	repositoryName := "longhorn"
	chartName := "longhorn"
	releaseName := "longhorn"

	chart := builder.NewChart(namespace)

	snapshotRetain := 1
	snapshotConcurrency := 1

	backupRetain := 3
	backupConcurrency := 1

	// longhorn-crypto
	bitnamicom.NewSealedSecret(
		chart.Cdk8sChart,
		jsii.String("longhorn-crypto-sealed-secret"),
		&bitnamicom.SealedSecretProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("longhorn-crypto"),
				Namespace: jsii.String(namespace),
			},
			Spec: &bitnamicom.SealedSecretSpec{
				Template: &bitnamicom.SealedSecretSpecTemplate{
					Metadata: map[string]string{
						"name":      "longhorn-crypto",
						"namespace": namespace,
					},
				},
				EncryptedData: &map[string]*string{
					"CRYPTO_KEY_CIPHER":   jsii.String(cryptKeyCipherValue),
					"CRYPTO_KEY_HASH":     jsii.String(cryptKeyHashValue),
					"CRYPTO_KEY_PROVIDER": jsii.String(cryptKeyProviderValue),
					"CRYPTO_KEY_SIZE":     jsii.String(cryptKeySizeValue),
					"CRYPTO_KEY_VALUE":    jsii.String(cryptoKeyValueValue),
					"CRYPTO_PBKDF_VALUE":  jsii.String(cryptoPBKDFValue),
				},
			},
		},
	)

	// storage class
	k8s.NewKubeStorageClass(
		chart.Cdk8sChart,
		jsii.String("longhorn-crypto-global"),
		&k8s.KubeStorageClassProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String("longhorn-crypto-global"),
				Annotations: &map[string]*string{
					"storageclass.kubernetes.io/is-default-class": jsii.String("true"),
				},
			},
			Provisioner:          jsii.String("driver.longhorn.io"),
			AllowVolumeExpansion: jsii.Bool(true),
			Parameters: &map[string]*string{
				"numberOfReplicas":    jsii.String("2"),
				"staleReplicaTimeout": jsii.String("2880"),
				"fromBackup":          jsii.String(""),
				"encrypted":           jsii.String("true"),
				"csi.storage.k8s.io/provisioner-secret-name":       jsii.String("longhorn-crypto"),
				"csi.storage.k8s.io/provisioner-secret-namespace":  jsii.String("longhorn-system"),
				"csi.storage.k8s.io/node-publish-secret-name":      jsii.String("longhorn-crypto"),
				"csi.storage.k8s.io/node-publish-secret-namespace": jsii.String("longhorn-system"),
				"csi.storage.k8s.io/node-stage-secret-name":        jsii.String("longhorn-crypto"),
				"csi.storage.k8s.io/node-stage-secret-namespace":   jsii.String("longhorn-system"),
			},
			ReclaimPolicy: jsii.String("Retain"),
		},
	)

	chart.NewNamespace(namespace)
	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)

	chart.CreateHelmRepository(
		repositoryName,
		"https://charts.longhorn.io",
	)

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "nas0-minio")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "basic-auth-users")

	certificates_certmanagerio.NewCertificate(
		chart.Cdk8sChart,
		jsii.String("certificate"),
		&certificates_certmanagerio.CertificateProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("secret-tls-www"),
			},
			Spec: &certificates_certmanagerio.CertificateSpec{
				DnsNames: &[]*string{
					jsii.String("longhorn.services.mkz.me"),
				},
				IssuerRef: &certificates_certmanagerio.CertificateSpecIssuerRef{
					Kind: jsii.String("ClusterIssuer"),
					Name: jsii.String("letsencrypt-prod"),
				},
				SecretName: jsii.String("secret-tls-www"),
			},
		},
	)

	// The following is no longer useful.
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
					Realm:  jsii.String("Longhorn Authentication"),
					Secret: jsii.String("basic-auth-users"),
				},
			},
		},
	)

	traefikio.NewIngressRoute(
		chart.Cdk8sChart,
		jsii.String("ingress-route"),
		&traefikio.IngressRouteProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
			},
			Spec: &traefikio.IngressRouteSpec{
				EntryPoints: &[]*string{
					// jsii.String("web"),
					jsii.String("websecure"),
				},
				Routes: &[]*traefikio.IngressRouteSpecRoutes{
					{
						Kind:  traefikio.IngressRouteSpecRoutesKind_RULE,
						Match: jsii.String("Host(`longhorn.services.mkz.me`)"),
						Middlewares: &[]*traefikio.IngressRouteSpecRoutesMiddlewares{
							// {
							// 	Name:      jsii.String("basic-auth"),
							// 	Namespace: jsii.String(namespace),
							// },
							{
								Name:      jsii.String("traefik-forward-auth"),
								Namespace: jsii.String("traefik-forward-auth"),
							},
						},
						Services: &[]*traefikio.IngressRouteSpecRoutesServices{
							{
								Kind: traefikio.IngressRouteSpecRoutesServicesKind_SERVICE,
								Name: jsii.String("longhorn-frontend"),
								Port: traefikio.IngressRouteSpecRoutesServicesPort_FromNumber(jsii.Number(80)),
							},
						},
					},
				},
				Tls: &traefikio.IngressRouteSpecTls{
					SecretName: jsii.String("secret-tls-www"),
				},
			},
		},
	)

	// PV backups
	longhornio.NewRecurringJobV1Beta2(
		chart.Cdk8sChart,
		jsii.String("longhorn-backups"),
		&longhornio.RecurringJobV1Beta2Props{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("longhorn-backups"),
			},
			Spec: &longhornio.RecurringJobV1Beta2Spec{
				Cron: jsii.String("45 6 * * *"),
				Groups: &[]*string{
					jsii.String("default"),
				},
				Retain:      jsii.Number(backupRetain),
				Concurrency: jsii.Number(backupConcurrency),
				Labels: &map[string]*string{
					"job": jsii.String("daily-backup"),
				},
				Task: longhornio.RecurringJobV1Beta2SpecTask_BACKUP,
			},
		},
	)

	longhornio.NewRecurringJobV1Beta2(
		chart.Cdk8sChart,
		jsii.String("longhorn-backups-disabled"),
		&longhornio.RecurringJobV1Beta2Props{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("longhorn-backups-disabled"),
			},
			Spec: &longhornio.RecurringJobV1Beta2Spec{
				Cron: jsii.String("0 5 31 2 *"),
				Groups: &[]*string{
					jsii.String("disabled"),
				},
				Retain:      jsii.Number(backupRetain),
				Concurrency: jsii.Number(backupConcurrency),
				Labels: &map[string]*string{
					"job": jsii.String("daily-backup"),
				},
				Task: longhornio.RecurringJobV1Beta2SpecTask_BACKUP,
			},
		},
	)

	longhornio.NewRecurringJobV1Beta2(
		chart.Cdk8sChart,
		jsii.String("longhorn-snapshots"),
		&longhornio.RecurringJobV1Beta2Props{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("longhorn-snapshots"),
			},
			Spec: &longhornio.RecurringJobV1Beta2Spec{
				Cron: jsii.String("15 1,7,13,19 * * *"),
				Groups: &[]*string{
					jsii.String("default"),
				},
				Retain:      jsii.Number(snapshotRetain),
				Concurrency: jsii.Number(snapshotConcurrency),
				Labels: &map[string]*string{
					"job": jsii.String("multiple-snapshot"),
				},
				Task: longhornio.RecurringJobV1Beta2SpecTask_SNAPSHOT,
			},
		},
	)

	longhornio.NewRecurringJobV1Beta2(
		chart.Cdk8sChart,
		jsii.String("longhorn-snapshots-disabled"),
		&longhornio.RecurringJobV1Beta2Props{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("longhorn-snapshots-disabled"),
			},
			Spec: &longhornio.RecurringJobV1Beta2Spec{
				Cron: jsii.String("0 5 31 2 *"),
				Groups: &[]*string{
					jsii.String("disabled"),
				},
				Retain:      jsii.Number(snapshotRetain),
				Concurrency: jsii.Number(snapshotConcurrency),
				Labels: &map[string]*string{
					"job": jsii.String("multiple-snapshot"),
				},
				Task: longhornio.RecurringJobV1Beta2SpecTask_SNAPSHOT,
			},
		},
	)

	longhornio.NewRecurringJobV1Beta2(
		chart.Cdk8sChart,
		jsii.String("longhorn-snapshots-daily"),
		&longhornio.RecurringJobV1Beta2Props{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("longhorn-snapshots-daily"),
			},
			Spec: &longhornio.RecurringJobV1Beta2Spec{
				Cron: jsii.String("0 5 31 2 *"),
				Groups: &[]*string{
					jsii.String("daily"),
				},
				Retain:      jsii.Number(snapshotRetain),
				Concurrency: jsii.Number(snapshotConcurrency),
				Labels: &map[string]*string{
					"job": jsii.String("multiple-daily"),
				},
				Task: longhornio.RecurringJobV1Beta2SpecTask_SNAPSHOT,
			},
		},
	)

	// Adding service monitor
	servicemonitor_monitoringcoreoscom.NewServiceMonitor(
		chart.Cdk8sChart,
		jsii.String("sm"),
		&servicemonitor_monitoringcoreoscom.ServiceMonitorProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Labels: &map[string]*string{
					"release": jsii.String("prometheus"),
				},
			},
			Spec: &servicemonitor_monitoringcoreoscom.ServiceMonitorSpec{
				Selector: &servicemonitor_monitoringcoreoscom.ServiceMonitorSpecSelector{
					MatchLabels: &map[string]*string{
						"app": jsii.String("longhorn-manager"),
					},
				},
				NamespaceSelector: &servicemonitor_monitoringcoreoscom.ServiceMonitorSpecNamespaceSelector{
					MatchNames: &[]*string{
						jsii.String("longhorn-system"),
					},
				},
				Endpoints: &[]*servicemonitor_monitoringcoreoscom.ServiceMonitorSpecEndpoints{
					{
						Port: jsii.String("manager"),
					},
				},
			},
		},
	)

	return chart
}

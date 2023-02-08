package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"

	"git.mkz.me/mycroft/k8s-home/imports/bitnamicom"
	"git.mkz.me/mycroft/k8s-home/imports/certificates_certmanagerio"
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/imports/traefikcontainous"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

const (
	CRYPTO_KEY_CIPHER_VALUE   = "AgBs1HRoXzlFCPX3WApqH1QxaJTJKzfrrJx82JxkEm0zT2LkOErfezIBIybstRYhIT0G5uYRVtwSCs9pxpOJIwofGs275Minpzcs8qsJCkZVU+1FQFemTGh4aevAOwHaV1SaimQy0tAWZsrRZ3TAbJfTGGbKanjsQ3ZprPfLBUoLgIjfXUdGY2mjwcVWjsCLWHj2Q/T2qGEGlB7yaK8NAmOY4G6tAvxT9rHKYG22Td1vKH0LzLUOjWkvF0KoxaEV1jdYjCtmV3GtiksjqyUgr4mP21svq3T5CeJo/LhtzEqm3jiZtpJP1dQmNrh4MCS2rq4+mH9KxWxIFz+zIxEXhZDxzBcx3gWNeKuzdhyBGJ6cNbL+2dcdm6BvKlCZva9wgyd57I58sS2ERWbOUS9L+VFIrdxJDAmu0x8Iy2v2zDGnIftj48U4zl+R2N5GUcpEFAJu8BWKC8AVYAI3z48Jq8o2gHNFlrTmKxX+R3V4Sl5o4k5NRwOubuRF9ulEdpSl+aeFUMt8YPY6Sydq7mSvziaiCSmHxPXAoN0J2+7Wcp85Y7J1fm9gZGcb/fNkEMBE4A2Z3xcJbJ+UaldSy+u1WfVTLn7VgKZdQJgmGa2taxaYyS4HNZhGPIoggV2X/QlOwLRUw3isOwj5a6cLWFx3P4cmV8b4puX+kvD4Qyn9QJWulnHwfwfkfLYm5kTHAPIBPcZT35w9fRsxSC/V9RgQpfc="
	CRYPTO_KEY_HASH_VALUE     = "AgBdywpwYo5RkYUF9zF65wiyxkCFuIJ83vuwt+9HctjvzJVdMEULJOPZl/VhxKyCbGBwnlySSoYEGwdQwM2fDYbuT/HSQiwy508qE6wI4wtMydgNzO3RRLieceCaVjXRG0B056Nu8T+FsenAOnFjlNnCQqh1UvoPPPzG8YaWLpk+NzbJrMyTPm2LYY3fg91//LUaZnzaS5qihRS8RLUCiEKbie0PArAI99FvzvV2WDb2q/f9tt9qoV6NpiLZt1KOpYXz63KmroewjauqLI5lgkjootyFi2sYRNwWyXSrR1q2iNdwVrVVqbeI4/gX6ATA7+xn1RpRoLVC0pcgLU8CahLEPkzOiOZUZvMxcGZK4TS6wUW4L9HiOHJ2eWZKp1LuU0OgdcnN515sdfcvE+PjI+TYQRo6YEhX8tV4l8nA8RhQWkLHqR9jLTNKif+ky7lVekx9kmv9m9Boj9WLcua2qE7E44On9IMP2YczJ+J5tZCbo89uLA/Xl0JGKeyMVC2rmczkB0j7wx5RHvY+JN4HGEI3KlnBnBO0mwjV97XVAVIMW/9TB8yNx0mWtlDiMbRzDASCI5Truo/JeEr6BDJkfN4Sxp6AntsQr4/jJYntLyK8JPULSGWsO6qeLatwO+Xr3/0t328DqzA2xSexeVjeE5letUcofRYOrtk3ffiZVCUBD98XVstsETYqyz5ww9h38tV4EX+Xd3A="
	CRYPTO_KEY_PROVIDER_VALUE = "AgCkfLOzp9Si++tup/FIVMq6dxLiiNPuXF4xgatBOzw43vDG/yQAvF/EOxWRQe1zYXXi1R7sNmL6elzSeQZQ4jbKNlveajFVK2wqISmevX1uNytrhC9udOh+bKC8mAxZrGitiSpSI2whRvFp52plf2+vNwyoCH8Klm1kA4MAn0Vkx4rp5h+F/yiIAXwcTS0feXjmi4F3Eop/j6kqz/oFXkxyjcD00+7DImqZkBzgKSHH1DdY6f1D0Jl/KZ2H/5ZlFt1oWBV0Sst2Scqz3Di1Egzsu7FHYSDS1v1xhsDnXuPwmYI4AePSjCIt3JIUHMrFTN5qJwrta0hQsAK4qRtxbEimF/ARrvPXOW9NxH1yFWR2EXJFEd8cLexNqnmsg0k16O1wcgVcysmmRnlJZ8cylwuZE1t0e2T3i/tLZ7iHVk869OgNzajtpE4D50q5z6WnZ9e/L1kyBE/iSZl6QixIiZC6WgqCE5I8m7g0JZR0GisAX6rp437ReTC/ZUmXRVm4+v6DiCi1VxQ58oOtsKbn1T+WZ9phGeBZL194cVuTOdeECErMZpHPZOHrC27+fGXuOXdVgRAJuVPlXpvpyj7Hu1QW0ltZ7SR94g4NWuw1UyRWOukU49zMcEGLKAJUNwkmn7xGl/u0oGqZJZ6NPtpaRPMR7iTrFewuyeXbrW2m/PzC0KKXHEKbnRqyYMQBC0lmG8uQvRflKOo="
	CRYPTO_KEY_SIZE_VALUE     = "AgCyVBbO5Bn9NpEGoqAKZIDZtmEw8DYbN+u87U8iU03z6me9XShrje4ks5UfA33TdL7bUdyIzrxsGQnESnCaNFw/aFv78UqppSA8bf5QD6ozxu8oxqU7cARiSfLLeKeq+eGSJi9bwtKVMq1v69YFO+qEeUNPGWsPewCJvHgr4BmLU2a4TZIrGDbWPhv4XW2AinNHvkNJUIeMUuh94r0nnU1RKw2G/38AocSrKhW2V+4e6rX9J8dryBc1kKYw0F/x9RwoX32wNjN/2GDfhPUNn3KmB7j+ES4g/QHPMLvFxJ4Y8RrcuunO1/mUFFdSVtD3LuOUoU/DCXL8hvOLnAqF60THhvJCzqrxLUMOyx1Mi7/reXST7YbBawYW0j/W+cK3vuhxvN4HIRCAqe3Ix1J5zWadUNKCzXo/zrioG+qNDg9oTjhhyNXMNjbMOvZ5ZPiD7Mbkuy+40tqF6L+r8CpatHCdfl2nBSeggYQXvUcYFbkL9tAivBvB2Lfym207RU0NdaUwRPSXSNdZsg+Cw4jy4VYRwtcREVrbmRBkwmYYiPZhWk5xXKF7NbH3zZE0kpJqOO2c6087DfNgPdT1Xi1fyO0Q6n/mhhncLPIYj1kJ2q1ZEyATTzHo5bfznzALP5TSJJFrW1PFbLxhfNrFegysGaYYlvjllU1Aon6yqBd7iSP9gUvjRccavebURWwXpXppmsBVhBg="
	CRYPTO_KEY_VALUE_VALUE    = "AgAHza4E9Pm2o9+NHjaolZrJ/E32LU6a6KcyBg+am6bME48Wc/gjSlZNCwQALNtTb5NdD51Nr14d01p0WKZ2962fzB434eL0bH1pimqQQ6zYuT36K9Eb6HSsF8fzP7yHMWqfWZmxBXsNc++L33HL60SYn1nZsybtDPr4AO2e0VKKd7jQFFPFaKZET677VjEUbGVS6R3Q0R4YCXCwc9q6HILGohsdYJmqYcn7bMcqUc9aTfchK6eNYgFKFFEI4ra51KU9gzgKHHKncIfROPeDq6TKmoiFgX1RZm70nTObYT2pG5Tna2IafgkBDkHOPgMmcypgW4RY0I3hN4O5i09CZrSyUYt1t9arxO8M9HGq1xTF3MxN448pfwkMnD5Qq0Quib1Il+jbOoUg5xKEJh84sSkDCqGvnCmlRufhLmf8oUCuRYLVkyJgqfQjFSvTVDyj3a6uMKuP3s4zxQTdYQRKrjrfLf6FrVu+6zlhL7Il+O74uSKBevmDVtM5BhLYiYhqq7DcuswiT3x4GUe5G/twT7YTCvmIswxT+vbIe3Z8Dd1G4qnGD91lpguAAM8fEFuCwrHfVw4M4Zn1IjbDYjB4xdyE4017LuaxTnxyyGvLYZ0Cc+NfQSHRDXYW6pzwKj6yfTgaCEBfitmK7KvJMSUrzaqacJngXIuwY23MdMUftFZKIKSN/TpgcjYGFR/IM+bp/6JzmteTd5EypSL8SiOJpLFIfL4nK7ETRZzeZA=="
	CRYPTO_PBKDF_VALUE        = "AgDTO42cqVxOGa7vO2A9ZyeFL/8818b6vpdUJFdr26lycb9ocQp0LbK/ubLnSZiontRKFD95QYf5Hy+iPqHVVw9io6VgbI3I20A5cRofFNIZQD0I4ht4bLTqxlyYTGCmkY2GST/o2JJ1LopcYKZBDmuOZkCH7JaPk4AX9N5JDu7Dhm0oOpFh8zaRvO8fkRrHu4jx4e9b/NPGf7czFZ55cT+S3tNtvCNeyspS0uX6CVuKXqXjf8c7VAJXWNgPXm6Ye0Yjp7+OtiXIun5OD0uVV7XDUufdjhXkd0z9iXhScBZZEyYruKFTqIcCvm+ZOCkAUmkmQioUtV5rh+tBSv6ZA5jT3A0hIrbna6nSeyTE5orayfN/yQ/I8QgkrAkPIewhGxkbqLX2YFIOnrd4IFztb8r4OGdCMjP6oU8Qv1O6s1bC2N/SbpBsL0DrpZGgboUfTwC9s21Qy3aH8nB70B1Pu1hBhWouxmGyEQ8Rt6z4XR6SmNbCwMV6rQGp/hLONISeeDsDtaiT1tiQtB8YVM+e+qngLPMI7RrnACetPzuv0Oekii6qVSlzk9F7372lpaak417pNjnled4dHK54De+O1/+Y3rGKTrwuinEDAzn1wMLJ3wCD3W3W9E2CG+M//vNd5YbgeKcEq8AJiM4abao9dx9mnBAT1R8sqF4lPqCr88OMDlG3w09HO6mrNWgZgknwp/600nOmrcDi"
)

func NewLonghornChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "longhorn-system"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	// longhorn-crypto
	bitnamicom.NewSealedSecret(
		chart,
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
					"CRYPTO_KEY_CIPHER":   jsii.String(CRYPTO_KEY_CIPHER_VALUE),
					"CRYPTO_KEY_HASH":     jsii.String(CRYPTO_KEY_HASH_VALUE),
					"CRYPTO_KEY_PROVIDER": jsii.String(CRYPTO_KEY_PROVIDER_VALUE),
					"CRYPTO_KEY_SIZE":     jsii.String(CRYPTO_KEY_SIZE_VALUE),
					"CRYPTO_KEY_VALUE":    jsii.String(CRYPTO_KEY_VALUE_VALUE),
					"CRYPTO_PBKDF_VALUE":  jsii.String(CRYPTO_PBKDF_VALUE),
				},
			},
		},
	)

	// storage class
	k8s.NewKubeStorageClass(
		chart,
		jsii.String("longhorn-crypto-global"),
		&k8s.KubeStorageClassProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String("longhorn-crypto-global"),
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
		},
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
		"1.4.0",
		nil,
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"longhorn.yaml",
			),
		},
		nil,
	)

	k8s_helpers.CreateExternalSecret(chart, namespace, "basic-auth-users")

	certificates_certmanagerio.NewCertificate(
		chart,
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

	traefikcontainous.NewMiddleware(
		chart,
		jsii.String("basic-auth"),
		&traefikcontainous.MiddlewareProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name:      jsii.String("basic-auth"),
				Namespace: jsii.String(namespace),
			},
			Spec: &traefikcontainous.MiddlewareSpec{
				BasicAuth: &traefikcontainous.MiddlewareSpecBasicAuth{
					Realm:  jsii.String("Longhorn Authentication"),
					Secret: jsii.String("basic-auth-users"),
				},
			},
		},
	)

	traefikcontainous.NewIngressRoute(
		chart,
		jsii.String("ingress-route"),
		&traefikcontainous.IngressRouteProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
			},
			Spec: &traefikcontainous.IngressRouteSpec{
				EntryPoints: &[]*string{
					jsii.String("web"),
					jsii.String("websecure"),
				},
				Routes: &[]*traefikcontainous.IngressRouteSpecRoutes{
					{
						Kind:  traefikcontainous.IngressRouteSpecRoutesKind_RULE,
						Match: jsii.String("Host(`longhorn.services.mkz.me`)"),
						Middlewares: &[]*traefikcontainous.IngressRouteSpecRoutesMiddlewares{
							{
								Name:      jsii.String("basic-auth"),
								Namespace: jsii.String(namespace),
							},
						},
						Services: &[]*traefikcontainous.IngressRouteSpecRoutesServices{
							{
								Kind: traefikcontainous.IngressRouteSpecRoutesServicesKind_SERVICE,
								Name: jsii.String("longhorn-frontend"),
								Port: traefikcontainous.IngressRouteSpecRoutesServicesPort_FromNumber(jsii.Number(80)),
							},
						},
					},
				},
				Tls: &traefikcontainous.IngressRouteSpecTls{
					SecretName: jsii.String("secret-tls-www"),
				},
			},
		},
	)

	return chart
}

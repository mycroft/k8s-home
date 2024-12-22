package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/certificates_certmanagerio"
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/imports/traefikio"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewBookstackChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "bookstack"
	appName := namespace
	appImage := builder.RegisterContainerImage("linuxserver/bookstack")
	appPort := 80
	appIngress := "bookstack.services.mkz.me"

	useLegacyIngress := true

	chart := builder.NewChart(namespace)

	kubehelpers.NewNamespace(chart.Cdk8sChart, namespace)
	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "mariadb")

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	env := []*k8s.EnvVar{
		{Name: jsii.String("APP_URL"), Value: jsii.String(fmt.Sprintf("https://%s", appIngress))},
		{Name: jsii.String("PUID"), Value: jsii.String("1000")},
		{Name: jsii.String("PGID"), Value: jsii.String("1000")},
		{Name: jsii.String("TZ"), Value: jsii.String("Etc/UTC")},
		{Name: jsii.String("DB_HOST"), Value: jsii.String("mariadb.mariadb")},
		{Name: jsii.String("DB_PORT"), Value: jsii.String("3306")},
		{Name: jsii.String("DB_USER"), Value: jsii.String("bookstack")},
		{Name: jsii.String("DB_DATABASE"), Value: jsii.String("bookstack")},
		{
			Name: jsii.String("DB_PASS"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("mariadb"),
					Key:  jsii.String("password"),
				},
			},
		},
		{Name: jsii.String("SESSION_LIFETIME"), Value: jsii.String("1800")},
	}

	_, serviceName := kubehelpers.NewStatefulSet(
		chart.Cdk8sChart,
		namespace,
		appName,
		appImage,
		appPort,
		labels,
		env,
		[]string{},
		[]kubehelpers.ConfigMapMount{},
		[]kubehelpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/config",
				StorageSize: "1Gi",
			},
		},
	)

	if useLegacyIngress {
		kubehelpers.NewAppIngress(
			builder.Context,
			chart.Cdk8sChart,
			labels,
			appName,
			appPort,
			appIngress,
			serviceName,
			map[string]string{},
		)

	} else {
		certificates_certmanagerio.NewCertificate(
			chart.Cdk8sChart,
			jsii.String("certificate"),
			&certificates_certmanagerio.CertificateProps{
				Metadata: &cdk8s.ApiObjectMetadata{
					Name:      jsii.String("ingress-tls-secret"),
					Namespace: jsii.String("bookstack"),
				},
				Spec: &certificates_certmanagerio.CertificateSpec{
					DnsNames: jsii.Strings("bookstack.services.mkz.me"),
					IssuerRef: &certificates_certmanagerio.CertificateSpecIssuerRef{
						Group: jsii.String("cert-manager.io"),
						Kind:  jsii.String("ClusterIssuer"),
						Name:  jsii.String("letsencrypt-prod"),
					},
					SecretName: jsii.String("ingress-tls-secret"),
					Usages: &[]certificates_certmanagerio.CertificateSpecUsages{
						certificates_certmanagerio.CertificateSpecUsages_DIGITAL_SIGNATURE,
						certificates_certmanagerio.CertificateSpecUsages_KEY_ENCIPHERMENT,
					},
				},
			},
		)

		traefikio.NewIngressRoute(
			chart.Cdk8sChart,
			jsii.Sprintf("ingress-route"),
			&traefikio.IngressRouteProps{
				Metadata: &cdk8s.ApiObjectMetadata{
					Name:      jsii.String("bookstack"),
					Namespace: jsii.String("bookstack"),
				},
				Spec: &traefikio.IngressRouteSpec{
					Routes: &[]*traefikio.IngressRouteSpecRoutes{
						{
							Kind:  traefikio.IngressRouteSpecRoutesKind_RULE,
							Match: jsii.String("Host(`bookstack.services.mkz.me`)"),
							Middlewares: &[]*traefikio.IngressRouteSpecRoutesMiddlewares{
								{
									Name:      jsii.String("ak-outpost-authentik-embedded-outpost"),
									Namespace: jsii.String("authentik"),
								},
							},
							Services: &[]*traefikio.IngressRouteSpecRoutesServices{
								{
									Name: jsii.String(serviceName),
									Port: traefikio.IngressRouteSpecRoutesServicesPort_FromString(jsii.String("http")),
								},
							},
						},
					},
					Tls: &traefikio.IngressRouteSpecTls{
						SecretName: jsii.String("ingress-tls-secret"),
					},
				},
			},
		)
	}

	return chart
}

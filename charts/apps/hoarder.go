package apps

import (
	"fmt"
	"strings"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

// https://docs.hoarder.app/Installation/docker
// https://raw.githubusercontent.com/hoarder-app/hoarder/main/docker/docker-compose.yml

func NewHoarderChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "hoarder"
	appName := namespace
	appImage := builder.RegisterContainerImage("ghcr.io/hoarder-app/hoarder")
	appPort := 3000
	appIngress := "hoarder.services.mkz.me"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)
	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	chart.CreateExternalSecret(namespace, "openai")
	chart.CreateExternalSecret(namespace, "sso")
	chart.CreateExternalSecret(namespace, "secret")

	labels := map[string]*string{
		"app.kubernetes.io/name":       jsii.String(appName),
		"app.kubernetes.io/components": jsii.String("hoarder"),
	}

	meiliLabels := map[string]*string{
		"app.kubernetes.io/name":       jsii.String(appName),
		"app.kubernetes.io/components": jsii.String("meili"),
	}

	chromiumLabels := map[string]*string{
		"app.kubernetes.io/name":       jsii.String(appName),
		"app.kubernetes.io/components": jsii.String("chromium"),
	}

	chromiumArgs := []string{
		"chromium-browser --headless",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"--disable-gpu",
		"--remote-debugging-address=0.0.0.0",
		"--remote-debugging-port=9222",
		"--hide-scrollbars",
	}

	kubehelpers.NewAppDeployment(
		chart.Cdk8sChart,
		namespace,
		"alpine-chrome",
		"gcr.io/zenika-hub/alpine-chrome:123",
		chromiumLabels,
		[]*k8s.EnvVar{},
		[]string{
			strings.Join(chromiumArgs, " "),
		},
		[]kubehelpers.ConfigMapMount{},
	)

	chromiumSvcName := kubehelpers.NewAppService(
		chart.Cdk8sChart,
		namespace,
		"chromium",
		chromiumLabels,
		"chromium",
		9222,
	)

	_, meiliSvcName := kubehelpers.NewStatefulSet(
		chart.Cdk8sChart,
		namespace,
		"meilisearch",
		"getmeili/meilisearch:v1.11.1",
		7700,
		meiliLabels,
		[]*k8s.EnvVar{
			{Name: jsii.String("MEILI_NO_ANALYTICS"), Value: jsii.String("true")},
		},
		[]string{},
		[]kubehelpers.ConfigMapMount{},
		[]kubehelpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/meili_data",
				StorageSize: "32Gi",
			},
		},
	)

	env := []*k8s.EnvVar{
		{Name: jsii.String("NEXTAUTH_URL"), Value: jsii.String(fmt.Sprintf("https://%s", appIngress))},
		{
			Name: jsii.String("NEXTAUTH_SECRET"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Key:  jsii.String("secret"),
					Name: jsii.String("secret"),
				},
			},
		},

		{Name: jsii.String("DATA_DIR"), Value: jsii.String("/data")},
		{Name: jsii.String("BROWSER_WEB_URL"), Value: jsii.String(fmt.Sprintf("http://%s:9222", *chromiumSvcName.Name()))},
		{Name: jsii.String("MAX_ASSET_SIZE_MB"), Value: jsii.String("16")},
		{Name: jsii.String("MEILI_ADDR"), Value: jsii.String(fmt.Sprintf("http://%s:7700", meiliSvcName))},
		{
			Name: jsii.String("OPENAI_API_KEY"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Key:  jsii.String("api_key"),
					Name: jsii.String("openai"),
				},
			},
		},
		{Name: jsii.String("DISABLE_PASSWORD_AUTH"), Value: jsii.String("true")},
		{
			Name: jsii.String("OAUTH_WELLKNOWN_URL"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Key:  jsii.String("discovery_url"),
					Name: jsii.String("sso"),
				},
			},
		},
		{
			Name: jsii.String("OAUTH_CLIENT_SECRET"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Key:  jsii.String("client_secret"),
					Name: jsii.String("sso"),
				},
			},
		},
		{
			Name: jsii.String("OAUTH_CLIENT_ID"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Key:  jsii.String("client_id"),
					Name: jsii.String("sso"),
				},
			},
		},
		{
			Name:  jsii.String("OAUTH_PROVIDER_NAME"),
			Value: jsii.String("Authentik"),
		},
	}

	_, svcName := kubehelpers.NewStatefulSet(
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
				MountPath:   "/data",
				StorageSize: "32Gi",
			},
		},
	)

	kubehelpers.NewAppIngress(
		builder.Context,
		chart.Cdk8sChart,
		labels,
		appName,
		appPort,
		appIngress,
		svcName,
		map[string]string{},
	)

	return chart
}

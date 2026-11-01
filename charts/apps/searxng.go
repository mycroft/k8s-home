package apps

import (
	"log"
	"os"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

// NewSearxngChart deploys SearXNG (https://docs.searxng.org), a metasearch
// engine, at https://searxng.merino-tiaki.ts.net on the tailnet only: served
// from the internet it would be an open search proxy, and the engines it
// queries throttle or block the IP behind one. In-cluster clients reach it at
// http://searxng.searxng.svc:8080, where format=json serves its API.
//
// configs/searxng/settings.yml is merged over the upstream defaults.
func NewSearxngChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "searxng"
	appName := namespace
	appImage := builder.RegisterContainerImage("searxng/searxng")
	appPort := uint(8080)

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	// Generated in-cluster, so there is no Vault entry to create.
	kubehelpers.CreateGeneratedPasswordSecret(chart.Cdk8sChart, namespace, "secret-key")

	labels := map[string]string{
		"app.kubernetes.io/name": appName,
	}

	settings, err := os.ReadFile("configs/searxng/settings.yml")
	if err != nil {
		log.Fatalf("Could not read searxng settings: %s", err.Error())
	}

	configMap := k8s.NewKubeConfigMap(
		chart.Cdk8sChart,
		jsii.String("config"),
		&k8s.KubeConfigMapProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
			},
			Data: &map[string]*string{
				"settings.yml": jsii.String(string(settings)),
			},
		},
	)

	env := []*k8s.EnvVar{
		{Name: jsii.String("SEARXNG_BASE_URL"), Value: jsii.String("https://" + appName + ".merino-tiaki.ts.net/")},
		{
			Name: jsii.String("SEARXNG_SECRET"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("secret-key"),
					Key:  jsii.String("password"),
				},
			},
		},
	}

	kubehelpers.NewAppDeployment(
		chart.Cdk8sChart,
		namespace,
		appName,
		appImage,
		kubehelpers.ToLabelsPtr(labels),
		env,
		[]string{},
		[]kubehelpers.ConfigMapMount{
			{
				Name:      "config",
				ConfigMap: configMap,
				MountPath: "/etc/searxng",
			},
		},
		kubehelpers.AppDeploymentOption{
			// The image declares no USER, so SearXNG would run as root. Its
			// files and /var/cache/searxng belong to searxng, uid 977. The
			// entrypoint warns that the mounted /etc/searxng is not owned by
			// searxng; it only reads it.
			RunAsNonRoot: true,
			RunAsUser:    977,
			PodAnnotations: map[string]string{
				"configMapHash": kubehelpers.ComputeConfigMapHash(configMap),
			},
		},
	)

	svc := kubehelpers.NewAppService(
		chart.Cdk8sChart,
		namespace,
		"svc",
		kubehelpers.ToLabelsPtr(labels),
		"http",
		appPort,
		// A fixed name, for the stable in-cluster URL above.
		kubehelpers.AppServiceOption{Name: appName},
	)

	chart.NewTailscaleIngress(appName, *svc.Name(), appPort)

	return chart
}

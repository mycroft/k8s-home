package apps

import (
	"log"
	"os"
	"strconv"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

// NewLiteLLMChart deploys the LiteLLM proxy (https://docs.litellm.ai/): one
// OpenAI-compatible endpoint in front of several model providers, choosing
// the backend by the requested model name. It routes to llama-server on
// moonstone and to OpenAI; configs/litellm/config.yaml lists the models.
//
// It serves https://litellm.services.mkz.me, reachable from the internet, and
// https://litellm.merino-tiaki.ts.net on the tailnet. Every request must carry
// the master key as bearer token. It runs without a database: no virtual keys,
// no spend tracking, and the admin UI is disabled.
//
// The secrets live in Vault at secret/namespaces/litellm/litellm:
//   - master-key: the clients' bearer token; LiteLLM requires the sk- prefix
//   - llama-api-key: llama-server's API key
//   - openai-api-key: OpenAI's API key
//
// The pod reads them at startup, so restart it after changing them.
//
// LiteLLM's release pipeline was compromised on 2026-10-21: v1.107.0 and the
// latest, main-latest and main-stable tags shipped a credential stealer.
// Pin release tags only, and review each bump before merging it.
func NewLiteLLMChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "litellm"
	appName := namespace
	appImage := builder.RegisterContainerImage("ghcr.io/berriai/litellm-non_root")
	appPort := uint(4000)
	appIngress := "litellm.services.mkz.me"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "litellm")

	const configFile = "config.yaml"
	configContents, err := os.ReadFile("configs/litellm/" + configFile)
	if err != nil {
		log.Fatalf("read litellm configuration: %v", err)
	}

	configMap := k8s.NewKubeConfigMap(
		chart.Cdk8sChart,
		jsii.String("config"),
		&k8s.KubeConfigMapProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
			},
			Data: &map[string]*string{
				configFile: jsii.String(string(configContents)),
			},
		},
	)

	labels := map[string]string{
		"app.kubernetes.io/name": appName,
	}

	secretEnv := func(name, key string) *k8s.EnvVar {
		return &k8s.EnvVar{
			Name: jsii.String(name),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("litellm"),
					Key:  jsii.String(key),
				},
			},
		}
	}

	env := []*k8s.EnvVar{
		// Do not load credentials from a .env file.
		{Name: jsii.String("LITELLM_MODE"), Value: jsii.String("PRODUCTION")},
		// Nothing here needs the admin UI or the API reference, and both would
		// be reachable from the internet: serve neither.
		{Name: jsii.String("DISABLE_ADMIN_UI"), Value: jsii.String("True")},
		{Name: jsii.String("NO_DOCS"), Value: jsii.String("True")},
		{Name: jsii.String("NO_REDOC"), Value: jsii.String("True")},
		// Read by the configuration's os.environ/ values.
		secretEnv("LITELLM_MASTER_KEY", "master-key"),
		secretEnv("LLAMA_API_KEY", "llama-api-key"),
		secretEnv("OPENAI_API_KEY", "openai-api-key"),
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
				MountPath: "/etc/litellm",
			},
		},
		kubehelpers.AppDeploymentOption{
			Args: []string{
				"--config", "/etc/litellm/" + configFile,
				"--port", strconv.FormatUint(uint64(appPort), 10),
				// One worker per pod, as LiteLLM recommends on Kubernetes.
				"--num_workers", "1",
			},
			// The image runs as uid 65534.
			RunAsNonRoot: true,
			// The configuration is only read at startup.
			PodAnnotations: map[string]string{
				"configMapHash": kubehelpers.ComputeConfigMapHash(configMap),
			},
		},
	)

	_, svcName := chart.NewIngress(&kubehelpers.Ingress{
		Name:      appName,
		Port:      appPort,
		Ingresses: []string{appIngress},
		Labels:    labels,
	})

	chart.NewTailscaleIngress(appName, svcName, appPort)

	return chart
}

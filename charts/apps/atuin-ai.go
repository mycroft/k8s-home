package apps

import (
	"log"
	"os"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

// NewAtuinAIChart deploys the Atuin AI server
// (https://github.com/atuinsh/atuin-ai-server), the backend of atuin ai: it
// turns a request typed in the shell into a suggested command, asking the
// model served by llama-server on moonstone. Its configuration, models
// included, is configs/atuin-ai/config.toml.
//
// It serves https://atuin-ai.services.mkz.me, reachable from the internet, and
// https://atuin-ai.merino-tiaki.ts.net on the tailnet. Every request must carry
// the AUTH_TOKEN bearer token: without one set, the server answers anyone.
// Both secrets live in Vault at secret/namespaces/atuin-ai/atuin-ai-server:
//   - auth-token: what clients set as api_token in their [ai] section
//   - llama-api-key: llama-server's API key
//
// The pod reads them at startup, so restart it after changing them.
func NewAtuinAIChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "atuin-ai"
	appName := namespace
	// Upstream publishes only the main and latest tags, so versions.yaml pins
	// a digest, which check-versions cannot bump: bump it by hand.
	appImage := builder.RegisterContainerImage("ghcr.io/atuinsh/atuin-ai-server")
	appPort := uint(8080)
	appIngress := "atuin-ai.services.mkz.me"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "atuin-ai-server")

	const configFile = "config.toml"
	configContents, err := os.ReadFile("configs/atuin-ai/" + configFile)
	if err != nil {
		log.Fatalf("read atuin-ai configuration: %v", err)
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

	env := []*k8s.EnvVar{
		{
			Name: jsii.String("AUTH_TOKEN"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("atuin-ai-server"),
					Key:  jsii.String("auth-token"),
				},
			},
		},
		{
			// Read by api_key in the configuration.
			Name: jsii.String("LLAMA_API_KEY"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("atuin-ai-server"),
					Key:  jsii.String("llama-api-key"),
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
				MountPath: "/etc/atuin-ai",
			},
		},
		kubehelpers.AppDeploymentOption{
			// The image runs as USER nobody, uid 65534.
			RunAsNonRoot: true,
			RunAsUser:    65534,
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

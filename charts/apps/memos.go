package apps

import (
	"fmt"
	"log"
	"os"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	kube "git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewMemosChart(builder *kube.Builder) *kube.Chart {
	name := "memos"
	namespace := name
	port := uint(5230)
	ingresses := []string{
		fmt.Sprintf("%s.services.mkz.me", name),
	}

	labels := map[string]string{
		"app.kubernetes.io/name": name,
	}

	env := []kube.EnvEntry{
		{Name: "MEMOS_DRIVER", Value: kube.EnvValue{Value: "postgres"}},
		{Name: "MEMOS_PORT", Value: kube.EnvValue{Value: "5230"}},
		{Name: "MEMOS_DSN", Value: kube.EnvValue{ValueFromSecret: kube.EnvValueFromSecret{
			Key:  "dsn",
			Name: "postgres",
		}}},
	}

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)
	chart.CreateSecretStore(namespace)
	chart.CreateExternalSecret(namespace, "postgres")

	const configFile = "memos-instance-setting-general.json"
	configContents, err := os.ReadFile("configs/memos/" + configFile)
	if err != nil {
		log.Fatalf("read memos configuration: %v", err)
	}

	k8s.NewKubeSecret(
		chart.Cdk8sChart,
		jsii.String("memos-config"),
		&k8s.KubeSecretProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("memos-config"),
				Namespace: jsii.String(namespace),
			},
			StringData: &map[string]*string{
				configFile: jsii.String(string(configContents)),
			},
		},
	)

	chart.NewDeployment(&kube.Deployment{
		Name:   name,
		Image:  "neosmemo/memos",
		Labels: labels,
		Env:    env,
		Secrets: []kube.SecretMount{
			{
				Name:      "memos-config",
				MountPath: "/etc/secrets",
			},
		},
	})

	chart.NewIngress(&kube.Ingress{
		Name:      name,
		Port:      port,
		Ingresses: ingresses,
		Labels:    labels,
	})

	return chart
}

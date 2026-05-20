package apps

import (
	"fmt"

	kube "git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
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

	chart.NewDeployment(&kube.Deployment{
		Name:   name,
		Image:  "neosmemo/memos",
		Labels: labels,
		Env:    env,
	})

	chart.NewIngress(&kube.Ingress{
		Name:      name,
		Port:      port,
		Ingresses: ingresses,
		Labels:    labels,
	})

	return chart
}

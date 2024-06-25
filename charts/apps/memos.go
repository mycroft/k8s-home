package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewMemosChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	name := "memos"
	namespace := name
	port := uint16(5230)
	ingresses := []string{
		fmt.Sprintf("%s.services.mkz.me", name),
	}

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(name),
	}

	env := []*k8s.EnvVar{
		{
			Name:  jsii.String("MEMOS_DRIVER"),
			Value: jsii.String("postgres"),
		},
		{
			Name: jsii.String("MEMOS_DSN"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Key:  jsii.String("dsn"),
					Name: jsii.String("postgres"),
				},
			},
		},
	}

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)
	chart.CreateSecretStore(namespace)
	chart.CreateExternalSecret(namespace, "postgres")

	chart.NewDeployment(&kubehelpers.Deployment{
		Namespace: namespace,
		Name:      name,
		Image:     "neosmemo/memos",
		Labels:    labels,
		Env:       env,
	})

	chart.NewIngress(&kubehelpers.Ingress{
		Namespace: namespace,
		Name:      name,
		Port:      port,
		Ingresses: ingresses,
		Labels:    labels,
	})

	return chart
}

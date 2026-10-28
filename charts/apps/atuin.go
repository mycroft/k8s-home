package apps

import (
	"strconv"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

// NewAtuinChart deploys the Atuin sync server (https://docs.atuin.sh), which
// keeps the shell history of every host running atuin. Records are encrypted
// on the client with a key the server never sees.
//
// It serves https://atuin.services.mkz.me, reachable from the internet, and
// https://atuin.merino-tiaki.ts.net on the tailnet. Its data lives in the
// atuin database on the CNPG cluster (charts/storage/cnpg-cluster.go), whose
// credentials the chart reads from secret/namespaces/atuin/postgresql.
func NewAtuinChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "atuin"
	appName := namespace
	appImage := builder.RegisterContainerImage("ghcr.io/atuinsh/atuin")
	appPort := uint(8888)
	appIngress := "atuin.services.mkz.me"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgresql")

	labels := map[string]string{
		"app.kubernetes.io/name": appName,
	}

	env := []*k8s.EnvVar{
		{Name: jsii.String("ATUIN_HOST"), Value: jsii.String("0.0.0.0")},
		{Name: jsii.String("ATUIN_PORT"), Value: jsii.String(strconv.FormatUint(uint64(appPort), 10))},
		// Open only until the account exists: one account serves every host,
		// the others join it with atuin login. Left open, anyone reaching the
		// public host can register.
		{Name: jsii.String("ATUIN_OPEN_REGISTRATION"), Value: jsii.String("true")},
		{
			Name: jsii.String("ATUIN_DB_URI"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("postgresql"),
					Key:  jsii.String("url"),
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
		[]kubehelpers.ConfigMapMount{},
		kubehelpers.AppDeploymentOption{
			// The entrypoint is atuin-server, which has no default command.
			Args: []string{"start"},
			// The image runs as USER atuin, uid 1000.
			RunAsNonRoot: true,
			RunAsUser:    1000,
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

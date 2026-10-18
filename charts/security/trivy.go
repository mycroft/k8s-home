package security

import (
	"log"
	"os"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewTrivyChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "trivy"
	namespace := "trivy-system"

	// Not "aqua": harbor-scanner-trivy already owns flux-system/aqua, pointing
	// at helm.aquasec.com, which does not carry trivy-operator.
	repoName := "aquasecurity"
	releaseName := "trivy-operator"

	chart := builder.NewChart(appName)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repoName,
		"https://aquasecurity.github.io/helm-charts/",
	)

	chart.CreateHelmRelease(
		namespace,
		repoName,
		"trivy-operator",
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
		// The cluster still carries CRDs from the 2023 install.
		kubehelpers.WithReplaceCRDs(),
	)

	// https://grafana.com/grafana/dashboards/17813, unmodified. Grafana's
	// dashboard sidecar loads labelled ConfigMaps from every namespace.
	dashboard, err := os.ReadFile("configs/trivy/grafana-dashboard.json")
	if err != nil {
		log.Fatalf("read trivy grafana dashboard: %v", err)
	}

	k8s.NewKubeConfigMap(
		chart.Cdk8sChart,
		jsii.String("grafana-dashboard"),
		&k8s.KubeConfigMapProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
				Labels: &map[string]*string{
					"grafana_dashboard": jsii.String("1"),
				},
			},
			Data: &map[string]*string{
				"trivy-operator.json": jsii.String(string(dashboard)),
			},
		},
	)

	return chart
}

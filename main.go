package main

import (
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	"git.mkz.me/mycroft/k8s-home/apps"
)

func main() {
	app := cdk8s.NewApp(nil)

	// security
	apps.NewCertManagerChart(app)
	apps.NewSealedSecretsChart(app)
	apps.NewVaultChart(app)
	apps.NewExternalSecretsChart(app)
	apps.NewTrivyChart(app)
	apps.NewDexIdpChart(app)
	apps.NewTraefikForwardAuth(app)

	// storage
	apps.NewLonghornChart(app)
	apps.NewPostgresOperator(app)
	apps.NewPostgres(app)
	apps.NewMinioOperator(app)
	apps.NewMinio(app)
	apps.NewScyllaOperatorChart(app)
	apps.NewScyllaChart(app)
	apps.NewNATSChart(app)
	apps.NewOpenSearchChart(app)
	apps.NewMariaDBOperator(app)
	apps.NewMariaDBChart(app)

	// observability
	apps.NewKubePrometheusStackChart(app)
	apps.NewLokiChart(app)
	apps.NewPromtailChart(app)

	// misc tooling
	apps.NewFluxCDChart(app)
	apps.NewKubernetesDashboardChart(app)

	// apps
	apps.NewHelloKubernetesChart(app)
	apps.NewWhatIsMyIpChart(app)
	apps.NewWallabagChart(app)
	apps.NewUrlsChart(app)
	apps.NewFreshRSS(app)
	apps.NewLinkdingChart(app)
	apps.NewPrivatebinChart(app)
	apps.NewPaperlessNGXChart(app)
	apps.NewYopassChart(app)
	apps.NewITToolsChart(app)

	app.Synth()
}

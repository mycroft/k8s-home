package main

import (
	"flag"

	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	"git.mkz.me/mycroft/k8s-home/apps"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
)

var checkVersion bool

func init() {
	flag.BoolVar(&checkVersion, "check-version", false, "check the installed versions")
}

func main() {
	flag.Parse()

	app := cdk8s.NewApp(nil)

	// security
	apps.NewCertManagerChart(app)
	apps.NewSealedSecretsChart(app)
	apps.NewVaultChart(app)
	apps.NewExternalSecretsChart(app)
	apps.NewTrivyChart(app)
	apps.NewDexIdpChart(app)
	apps.NewTraefikForwardAuth(app)
	apps.NewKyvernoChart(app)

	// storage
	apps.NewLonghornChart(app)
	apps.NewPostgresOperator(app)
	apps.NewPostgres(app)
	apps.NewMinioOperator(app)
	apps.NewMinio(app)
	apps.NewScyllaOperatorChart(app)
	apps.NewScyllaChart(app)
	apps.NewNATSChart(app)
	// apps.NewOpenSearchChart(app)
	apps.NewMariaDBOperator(app)
	apps.NewMariaDBChart(app)
	apps.NewVeleroChart(app)

	// observability
	apps.NewKubePrometheusStackChart(app)
	// apps.NewLokiChart(app)
	// apps.NewPromtailChart(app)
	// apps.NewJaegerChart(app)

	// misc tooling
	apps.NewFluxCDChart(app)
	apps.NewKubernetesDashboardChart(app)
	apps.NewLinkerdChart(app)

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
	apps.NewBookstackChart(app)
	apps.NewHeimdallChart(app)
	apps.NewEmojivotoChart(app)
	apps.NewVaultWardenChart(app)
	apps.NewSendChart(app)
	apps.NewHeyChart(app)
	apps.NewHappyUrlsChart(app)

	if checkVersion {
		k8s_helpers.CheckVersions()
	}

	app.Synth()
}

package main

import (
	"flag"

	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	"git.mkz.me/mycroft/k8s-home/apps"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"

	charts_apps "git.mkz.me/mycroft/k8s-home/charts/apps"
	charts_infra "git.mkz.me/mycroft/k8s-home/charts/infra"
	charts_storage "git.mkz.me/mycroft/k8s-home/charts/storage"
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
	charts_storage.NewLonghornChart(app)
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
	charts_infra.NewFluxCDChart(app)
	apps.NewKubernetesDashboardChart(app)
	apps.NewLinkerdChart(app)

	// apps
	charts_apps.NewHelloKubernetesChart(app)
	charts_apps.NewWhatIsMyIpChart(app)
	charts_apps.NewWallabagChart(app)
	charts_apps.NewUrlsChart(app)
	charts_apps.NewFreshRSS(app)
	charts_apps.NewLinkdingChart(app)
	charts_apps.NewPrivatebinChart(app)
	charts_apps.NewPaperlessNGXChart(app)
	charts_apps.NewYopassChart(app)
	charts_apps.NewITToolsChart(app)
	charts_apps.NewBookstackChart(app)
	charts_apps.NewHeimdallChart(app)
	charts_apps.NewEmojivotoChart(app)
	charts_apps.NewVaultWardenChart(app)
	charts_apps.NewSendChart(app)
	charts_apps.NewHeyChart(app)
	charts_apps.NewHappyUrlsChart(app)

	charts_apps.NewSnippetBoxChart(app)

	if checkVersion {
		k8s_helpers.CheckVersions()
	}

	app.Synth()
}

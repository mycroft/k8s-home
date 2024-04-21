package main

import (
	"flag"

	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"

	charts_apps "git.mkz.me/mycroft/k8s-home/charts/apps"
	charts_infra "git.mkz.me/mycroft/k8s-home/charts/infra"
	charts_observability "git.mkz.me/mycroft/k8s-home/charts/observability"
	charts_security "git.mkz.me/mycroft/k8s-home/charts/security"
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
	charts_security.NewCertManagerChart(app)
	charts_security.NewSealedSecretsChart(app)
	charts_security.NewVaultChart(app)
	charts_security.NewExternalSecretsChart(app)
	// charts_security.NewTrivyChart(app)
	charts_security.NewDexIdpChart(app)
	charts_security.NewTraefikForwardAuth(app)
	// charts_security.NewKyvernoChart(app)

	// storage
	charts_storage.NewLonghornChart(app)
	charts_storage.NewPostgresOperator(app)
	charts_storage.NewPostgres(app)
	charts_storage.NewMinioOperator(app)
	charts_storage.NewMinio(app)
	charts_storage.NewScyllaOperatorChart(app)
	charts_storage.NewScyllaChart(app)
	charts_storage.NewNATSChart(app)
	// apps.NewOpenSearchChart(app)
	charts_storage.NewMariaDBOperator(app)
	charts_storage.NewMariaDBChart(app)
	charts_infra.NewVeleroChart(app)

	// observability
	charts_observability.NewKubePrometheusStackChart(app)
	// charts_observability.NewLokiChart(app)
	// charts_observability.NewPromtailChart(app)
	// charts_observability.NewJaegerChart(app)

	// misc tooling
	charts_infra.NewFluxCDChart(app)
	charts_infra.NewKubernetesDashboardChart(app)
	// charts_infra.NewLinkerdChart(app)
	// charts_infra.NewTektonChart(app)

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
	// charts_apps.NewHeimdallChart(app)
	// charts_apps.NewEmojivotoChart(app)
	charts_apps.NewVaultWardenChart(app)
	charts_apps.NewSendChart(app)
	charts_apps.NewHeyChart(app)
	charts_apps.NewHappyUrlsChart(app)
	charts_apps.NewSnippetBoxChart(app)
	charts_apps.NewExcalidrawChart(app)
	// charts_apps.NewJitsiChart(app)
	charts_apps.NewWikiJsChart(app)
	charts_apps.NewRedmineChart(app)
	charts_apps.NewMicrobinChart(app)

	if checkVersion {
		k8s_helpers.CheckVersions()
	}

	app.Synth()
}

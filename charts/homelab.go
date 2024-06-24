package charts

import (
	"context"

	charts_apps "git.mkz.me/mycroft/k8s-home/charts/apps"
	charts_cicd "git.mkz.me/mycroft/k8s-home/charts/cicd"
	charts_infra "git.mkz.me/mycroft/k8s-home/charts/infra"
	charts_observability "git.mkz.me/mycroft/k8s-home/charts/observability"
	charts_security "git.mkz.me/mycroft/k8s-home/charts/security"
	charts_storage "git.mkz.me/mycroft/k8s-home/charts/storage"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func HomelabBuildApp(ctx context.Context) *kubehelpers.Builder {
	builder := kubehelpers.NewBuilder(ctx)

	charts := [](func(*kubehelpers.Builder) *kubehelpers.Chart){
		// security charts
		charts_security.NewCertManagerChart,
		charts_security.NewSealedSecretsChart,
		charts_security.NewVaultChart,
		charts_security.NewExternalSecretsChart,
		charts_security.NewDexIdpChart,
		charts_security.NewTraefikForwardAuth,
		charts_security.NewAuthentikChart,
		// charts_security.NewTrivyChart
		// charts_security.NewKyvernoChart

		// storage charts
		charts_storage.NewLonghornChart,
		charts_storage.NewPostgresOperator,
		charts_storage.NewPostgres,
		charts_storage.NewMinioOperator,
		charts_storage.NewMinio,
		charts_storage.NewScyllaOperatorChart,
		charts_storage.NewScyllaChart,
		charts_storage.NewNATSChart,
		charts_storage.NewMariaDBOperator,
		charts_storage.NewMariaDBChart,
		// charts_storage.NewOpenSearchChart,

		// misc infra charts
		charts_infra.NewFluxCDChart,
		charts_infra.NewCapacitorChart,
		charts_infra.NewVeleroChart,
		charts_infra.NewKubernetesDashboardChart,
		charts_infra.NewLinkerdChart,
		charts_infra.NewTektonChart,
		charts_infra.NewTemporalChart,
		charts_infra.NewTraefikChart,

		// observability
		charts_observability.NewGrafanaHelmRepositoryChart,
		charts_observability.NewKubePrometheusStackChart,
		charts_observability.NewBlackboxExporterChart,
		charts_observability.NewLokiChart,
		charts_observability.NewTempoChart,
		charts_observability.NewKarmaChart,
		charts_observability.NewAlloyChart,
		// charts_observability.NewPromtailChart
		// charts_observability.NewJaegerChart

		// apps
		charts_apps.NewBookstackChart,
		charts_apps.NewCalibreWebChart,
		charts_apps.NewCyberchefChart,
		charts_apps.NewEmojivotoChart,
		charts_apps.NewExcalidrawChart,
		charts_apps.NewFreshRSS,
		charts_apps.NewHappyUrlsChart,
		charts_apps.NewHelloKubernetesChart,
		charts_apps.NewHeyChart,
		charts_apps.NewHomepageChart,
		charts_apps.NewITToolsChart,
		charts_apps.NewLinkdingChart,
		charts_apps.NewMemosChart,
		charts_apps.NewMicrobinChart,
		charts_apps.NewPaperlessNGXChart,
		charts_apps.NewPrivatebinChart,
		charts_apps.NewRedmineChart,
		charts_apps.NewSendChart,
		charts_apps.NewSnippetBoxChart,
		charts_apps.NewUrlsChart,
		charts_apps.NewVaultWardenChart,
		charts_apps.NewWallabagChart,
		charts_apps.NewWhatIsMyIPChart,
		charts_apps.NewWikiJsChart,
		charts_apps.NewYopassChart,
		// charts_apps.NewHeimdallChart)
		// charts_apps.NewJitsiChart)

		// my own apps
		charts_apps.NewGossipChart,

		// misc: testing CI/CD
		charts_cicd.NewCICDChart,
	}

	for _, chartCallback := range charts {
		builder.BuildChart(chartCallback)
	}

	return builder
}

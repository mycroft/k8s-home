package charts

import (
	"context"

	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	charts_apps "git.mkz.me/mycroft/k8s-home/charts/apps"
	charts_cicd "git.mkz.me/mycroft/k8s-home/charts/cicd"
	charts_infra "git.mkz.me/mycroft/k8s-home/charts/infra"
	charts_observability "git.mkz.me/mycroft/k8s-home/charts/observability"
	charts_security "git.mkz.me/mycroft/k8s-home/charts/security"
	charts_storage "git.mkz.me/mycroft/k8s-home/charts/storage"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func HomelabBuildApp(ctx context.Context) cdk8s.App {
	builder := kubehelpers.NewBuilder(ctx)

	// security
	builder.BuildChart(charts_security.NewCertManagerChart)
	builder.BuildChart(charts_security.NewSealedSecretsChart)
	builder.BuildChartLegacy(charts_security.NewVaultChart)
	builder.BuildChartLegacy(charts_security.NewExternalSecretsChart)
	// charts_security.NewTrivyChart)
	builder.BuildChartLegacy(charts_security.NewDexIdpChart)
	builder.BuildChartLegacy(charts_security.NewTraefikForwardAuth)
	// charts_security.NewKyvernoChart)
	builder.BuildChartLegacy(charts_security.NewAuthentikChart)

	// storage
	builder.BuildChartLegacy(charts_storage.NewLonghornChart)
	builder.BuildChartLegacy(charts_storage.NewPostgresOperator)
	builder.BuildChartLegacy(charts_storage.NewPostgres)
	builder.BuildChartLegacy(charts_storage.NewMinioOperator)
	builder.BuildChartLegacy(charts_storage.NewMinio)
	builder.BuildChartLegacy(charts_storage.NewScyllaOperatorChart)
	builder.BuildChartLegacy(charts_storage.NewScyllaChart)
	builder.BuildChartLegacy(charts_storage.NewNATSChart)
	// apps.NewOpenSearchChart)
	builder.BuildChartLegacy(charts_storage.NewMariaDBOperator)
	builder.BuildChartLegacy(charts_storage.NewMariaDBChart)
	builder.BuildChartLegacy(charts_infra.NewVeleroChart)

	// observability
	builder.BuildChartLegacy(charts_observability.NewGrafanaHelmRepositoryChart)
	builder.BuildChartLegacy(charts_observability.NewKubePrometheusStackChart)
	builder.BuildChartLegacy(charts_observability.NewBlackboxExporterChart)
	builder.BuildChartLegacy(charts_observability.NewLokiChart)
	// charts_observability.NewPromtailChart)
	// charts_observability.NewJaegerChart)
	builder.BuildChartLegacy(charts_observability.NewTempoChart)
	builder.BuildChartLegacy(charts_observability.NewKarmaChart)
	builder.BuildChartLegacy(charts_observability.NewAlloyChart)

	// misc tooling
	builder.BuildChartLegacy(charts_infra.NewFluxCDChart)
	builder.BuildChartLegacy(charts_infra.NewCapacitorChart)
	builder.BuildChartLegacy(charts_infra.NewKubernetesDashboardChart)
	builder.BuildChartLegacy(charts_infra.NewLinkerdChart)
	builder.BuildChartLegacy(charts_infra.NewTektonChart)
	builder.BuildChartLegacy(charts_infra.NewTemporalChart)
	builder.BuildChartLegacy(charts_infra.NewTraefikChart)

	// apps
	builder.BuildChartLegacy(charts_apps.NewHelloKubernetesChart)
	builder.BuildChartLegacy(charts_apps.NewWhatIsMyIPChart)
	builder.BuildChartLegacy(charts_apps.NewWallabagChart)
	builder.BuildChartLegacy(charts_apps.NewUrlsChart)
	builder.BuildChartLegacy(charts_apps.NewFreshRSS)
	builder.BuildChartLegacy(charts_apps.NewLinkdingChart)
	builder.BuildChartLegacy(charts_apps.NewPrivatebinChart)
	builder.BuildChartLegacy(charts_apps.NewPaperlessNGXChart)
	builder.BuildChartLegacy(charts_apps.NewYopassChart)
	builder.BuildChartLegacy(charts_apps.NewITToolsChart)
	builder.BuildChartLegacy(charts_apps.NewBookstackChart)
	// charts_apps.NewHeimdallChart)
	builder.BuildChartLegacy(charts_apps.NewEmojivotoChart)
	builder.BuildChartLegacy(charts_apps.NewVaultWardenChart)
	builder.BuildChartLegacy(charts_apps.NewSendChart)
	builder.BuildChartLegacy(charts_apps.NewHeyChart)
	builder.BuildChartLegacy(charts_apps.NewHappyUrlsChart)
	builder.BuildChartLegacy(charts_apps.NewSnippetBoxChart)
	builder.BuildChartLegacy(charts_apps.NewExcalidrawChart)
	// charts_apps.NewJitsiChart)
	builder.BuildChartLegacy(charts_apps.NewWikiJsChart)
	builder.BuildChartLegacy(charts_apps.NewRedmineChart)
	builder.BuildChartLegacy(charts_apps.NewMicrobinChart)
	builder.BuildChartLegacy(charts_apps.NewCalibreWebChart)
	builder.BuildChartLegacy(charts_apps.NewHomepageChart)
	builder.BuildChartLegacy(charts_apps.NewCyberchefChart)

	// CI/CD
	builder.BuildChartLegacy(charts_cicd.NewCICDChart)

	return builder.App
}

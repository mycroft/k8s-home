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
)

func HomelabBuildApp(ctx context.Context) cdk8s.App {
	app := cdk8s.NewApp(nil)

	// security
	charts_security.NewCertManagerChart(ctx, app)
	charts_security.NewSealedSecretsChart(ctx, app)
	charts_security.NewVaultChart(ctx, app)
	charts_security.NewExternalSecretsChart(ctx, app)
	// charts_security.NewTrivyChart(ctx, app)
	charts_security.NewDexIdpChart(ctx, app)
	charts_security.NewTraefikForwardAuth(ctx, app)
	// charts_security.NewKyvernoChart(ctx, app)
	charts_security.NewAuthentikChart(ctx, app)

	// storage
	charts_storage.NewLonghornChart(ctx, app)
	charts_storage.NewPostgresOperator(ctx, app)
	charts_storage.NewPostgres(ctx, app)
	charts_storage.NewMinioOperator(ctx, app)
	charts_storage.NewMinio(ctx, app)
	charts_storage.NewScyllaOperatorChart(ctx, app)
	charts_storage.NewScyllaChart(ctx, app)
	charts_storage.NewNATSChart(ctx, app)
	// apps.NewOpenSearchChart(ctx, app)
	charts_storage.NewMariaDBOperator(ctx, app)
	charts_storage.NewMariaDBChart(ctx, app)
	charts_infra.NewVeleroChart(ctx, app)

	// observability
	charts_observability.NewGrafanaHelmRepositoryChart(ctx, app)
	charts_observability.NewKubePrometheusStackChart(ctx, app)
	charts_observability.NewBlackboxExporterChart(ctx, app)
	charts_observability.NewLokiChart(ctx, app)
	charts_observability.NewPromtailChart(ctx, app)
	// charts_observability.NewJaegerChart(ctx, app)
	charts_observability.NewTempoChart(ctx, app)

	// misc tooling
	charts_infra.NewFluxCDChart(ctx, app)
	charts_infra.NewCapacitorChart(ctx, app)
	charts_infra.NewKubernetesDashboardChart(ctx, app)
	charts_infra.NewLinkerdChart(ctx, app)
	charts_infra.NewTektonChart(ctx, app)
	charts_infra.NewTemporalChart(ctx, app)

	// apps
	charts_apps.NewHelloKubernetesChart(ctx, app)
	charts_apps.NewWhatIsMyIPChart(ctx, app)
	charts_apps.NewWallabagChart(ctx, app)
	charts_apps.NewUrlsChart(ctx, app)
	charts_apps.NewFreshRSS(ctx, app)
	charts_apps.NewLinkdingChart(ctx, app)
	charts_apps.NewPrivatebinChart(ctx, app)
	charts_apps.NewPaperlessNGXChart(ctx, app)
	charts_apps.NewYopassChart(ctx, app)
	charts_apps.NewITToolsChart(ctx, app)
	charts_apps.NewBookstackChart(ctx, app)
	// charts_apps.NewHeimdallChart(ctx, app)
	charts_apps.NewEmojivotoChart(ctx, app)
	charts_apps.NewVaultWardenChart(ctx, app)
	charts_apps.NewSendChart(ctx, app)
	charts_apps.NewHeyChart(ctx, app)
	charts_apps.NewHappyUrlsChart(ctx, app)
	charts_apps.NewSnippetBoxChart(ctx, app)
	charts_apps.NewExcalidrawChart(ctx, app)
	// charts_apps.NewJitsiChart(ctx, app)
	charts_apps.NewWikiJsChart(ctx, app)
	charts_apps.NewRedmineChart(ctx, app)
	charts_apps.NewMicrobinChart(ctx, app)

	// CI/CD
	charts_cicd.NewCICDChart(ctx, app)

	return app
}

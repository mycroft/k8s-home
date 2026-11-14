package apps

import (
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

// NewRssToBlueskyChart deploys rss-to-bluesky, which posts the newest
// lobste.rs stories to https://bsky.app/profile/lobsters-feed.bsky.social
// (https://github.com/mycroft/rss-to-bluesky), from its Helm chart: a CronJob
// with its state database on a PVC. Values are in configs/rss-to-bluesky.yaml.
//
// The Bluesky credentials live in Vault at
// secret/namespaces/rss-to-bluesky/rss-to-bluesky: bluesky-user and
// bluesky-password. Each run reads them at start.
//
// The database records the items already posted: an empty one makes the next
// run post the whole feed again. Keep it when moving the app, see "Moving an
// existing database" in the app's README.
func NewRssToBlueskyChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "rss-to-bluesky"
	repositoryName := "mycroft-ghcr"
	chartName := "rss-to-bluesky"
	releaseName := "rss-to-bluesky"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	chart.CreateHelmRepository(
		repositoryName,
		"oci://ghcr.io/mycroft/charts",
	)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "rss-to-bluesky")

	chart.CreateHelmRelease(
		namespace,
		repositoryName,
		chartName,
		releaseName,
		kubehelpers.WithDefaultConfigFile(),
	)

	return chart
}

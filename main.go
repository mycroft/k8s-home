package main

import (
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	"git.mkz.me/mycroft/k8s-home/apps"
)

func main() {
	app := cdk8s.NewApp(nil)
	apps.NewCertManagerChart(app)
	apps.NewKubePrometheusStackChart(app)
	apps.NewSealedSecretsChart(app)

	apps.NewHelloKubernetesChart(app)
	apps.NewWhatIsMyIpChart(app)

	app.Synth()
}

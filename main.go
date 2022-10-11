package main

import (
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	"git.mkz.me/mycroft/k8s-home/apps"
)

func main() {
	app := cdk8s.NewApp(nil)
	apps.NewHelloKubernetesChart(app)
	apps.NewCertManagerChart(app)
	app.Synth()
}

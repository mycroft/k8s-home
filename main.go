package main

import (
	"flag"

	"git.mkz.me/mycroft/k8s-home/charts"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
)

var checkVersion bool

func init() {
	flag.BoolVar(&checkVersion, "check-version", false, "check the installed versions")
}

func main() {
	flag.Parse()

	app := charts.HomelabBuildApp()

	if checkVersion {
		k8s_helpers.CheckVersions()
	}

	app.Synth()
}

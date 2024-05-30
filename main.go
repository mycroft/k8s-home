package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"git.mkz.me/mycroft/k8s-home/charts"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

var (
	debug  *bool
	filter *string
)

var rootCmd = &cobra.Command{
	Use:   "k8s-home",
	Short: "k8s-home is the yaml charts generator for my homelab",
	Run: func(cmd *cobra.Command, args []string) {
		GenerateYamlCharts()
	},
}

var GenerateYamlChartsCmd = &cobra.Command{
	Use:   "generate-yaml-charts",
	Short: "generates Yaml charts",
	Run: func(cmd *cobra.Command, args []string) {
		GenerateYamlCharts()
	},
}

var checkVersionCmd = &cobra.Command{
	Use:   "check-versions",
	Short: "check versions of declared helm charts & docker images",
	Run: func(cmd *cobra.Command, args []string) {
		if *debug {
			log.Println("preparing charts...")
		}
		_ = charts.HomelabBuildApp()

		if *debug {
			log.Println("running check-versions...")
		}
		kubehelpers.CheckVersions(*debug, *filter)
	},
}

func GenerateYamlCharts() {
	if *debug {
		log.Println("preparing charts...")
	}

	app := charts.HomelabBuildApp()

	if *debug {
		log.Println("syntheizing yamls...")
	}

	app.Synth()
}

func init() {
	rootCmd.AddCommand(checkVersionCmd)
	rootCmd.AddCommand(GenerateYamlChartsCmd)

	debug = rootCmd.PersistentFlags().Bool("debug", false, "enable debug")
	filter = checkVersionCmd.Flags().String("filter", "", "filter to apply when checking helm/container images")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package main

import (
	"context"
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
	Run: func(command *cobra.Command, _ []string) {
		GenerateYamlCharts(command.Flags().Lookup("versions").Value.String())
	},
}

var GenerateYamlChartsCmd = &cobra.Command{
	Use:     "generate-yaml-charts",
	Short:   "generates Yaml charts",
	Long:    "generates Yaml charts for all the helm charts used in my homelab",
	Example: "k8s-home generate-yaml-charts",
	Aliases: []string{"generate-yaml"},
	Run: func(command *cobra.Command, _ []string) {
		GenerateYamlCharts(command.Flags().Lookup("versions").Value.String())
	},
}

var checkVersionCmd = &cobra.Command{
	Use:   "check-versions",
	Short: "check versions of declared helm charts & docker images",
	Run: func(command *cobra.Command, _ []string) {
		if *debug {
			log.Println("preparing charts...")
		}

		builder := charts.HomelabBuildApp(
			context.WithValue(
				context.TODO(),
				kubehelpers.ValueKey,
				kubehelpers.ContextValues{
					Debug: *debug,
				},
			),
			command.Flags().Lookup("versions").Value.String(),
		)

		if *debug {
			log.Println("running check-versions...")
		}
		builder.CheckVersions(*debug, *filter)
	},
}

func GenerateYamlCharts(versionsFile string) {
	if *debug {
		log.Println("preparing charts...")
	}

	builder := charts.HomelabBuildApp(
		context.WithValue(context.TODO(),
			kubehelpers.ValueKey,
			kubehelpers.ContextValues{
				Debug: *debug,
			},
		),
		versionsFile,
	)

	if *debug {
		log.Println("syntheizing yamls...")
	}

	builder.App.Synth()
}

func init() {
	rootCmd.AddCommand(checkVersionCmd)
	rootCmd.AddCommand(GenerateYamlChartsCmd)

	rootCmd.PersistentFlags().String("versions", "versions.yaml", "versions.yaml file to user")

	debug = rootCmd.PersistentFlags().Bool("debug", false, "enable debug")
	filter = checkVersionCmd.Flags().String("filter", "", "filter to apply when checking helm/container images")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

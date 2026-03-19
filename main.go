package main

import (
	"context"
	"fmt"
	"log"
	"maps"
	"os"

	"github.com/spf13/cobra"

	"git.mkz.me/mycroft/k8s-home/charts"
	"git.mkz.me/mycroft/k8s-home/internal/gitea"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

var (
	debug            *bool
	filter           *string
	giteaURL         *string
	owner            *string
	repo             *string
	baseBranch       *string
	dryRun           *bool
	prFilter         *string
	prBranchOverride *string
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

var createPRsCmd = &cobra.Command{
	Use:   "create-prs",
	Short: "create pull requests for all outdated helm charts and container images",
	Run: func(command *cobra.Command, _ []string) {
		runCreatePRs(command)
	},
}

func runCreatePRs(command *cobra.Command) {
	ctx := context.Background()

	token := os.Getenv("GITEA_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	if token == "" {
		log.Fatal("GITEA_TOKEN or GITHUB_TOKEN environment variable must be set")
	}

	builder := charts.HomelabBuildApp(
		context.WithValue(ctx, kubehelpers.ValueKey, kubehelpers.ContextValues{Debug: *debug}),
		command.Flags().Lookup("versions").Value.String(),
	)

	helmUpdates, err := builder.GetHelmUpdates(*debug, *prFilter)
	if err != nil {
		log.Fatalf("get helm updates: %v", err)
	}

	imageUpdates, err := builder.GetImageUpdates(*debug, *prFilter)
	if err != nil {
		log.Fatalf("get image updates: %v", err)
	}

	allUpdates := make(map[string]string, len(helmUpdates)+len(imageUpdates))
	maps.Copy(allUpdates, helmUpdates)
	maps.Copy(allUpdates, imageUpdates)

	client := gitea.NewClient(*giteaURL, token, *owner, *repo)

	if err = gitea.CreateVersionBumpPRs(ctx, client, allUpdates, *baseBranch, *prBranchOverride, *dryRun); err != nil {
		log.Fatalf("create PRs: %v", err)
	}
}

var mergePRCmd = &cobra.Command{
	Use:   "merge-pr <number>",
	Short: "merge an open pull request using rebase",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		ctx := context.Background()

		token := os.Getenv("GITEA_TOKEN")
		if token == "" {
			token = os.Getenv("GITHUB_TOKEN")
		}

		if token == "" {
			log.Fatal("GITEA_TOKEN or GITHUB_TOKEN environment variable must be set")
		}

		var number int
		if _, err := fmt.Sscan(args[0], &number); err != nil {
			log.Fatalf("invalid PR number %q: %v", args[0], err)
		}

		client := gitea.NewClient(*giteaURL, token, *owner, *repo)

		if err := client.MergePR(ctx, number); err != nil {
			log.Fatalf("merge PR: %v", err)
		}

		fmt.Printf("PR #%d merged.\n", number)
	},
}

var listPRsCmd = &cobra.Command{
	Use:   "list-prs",
	Short: "list open pull requests in the repository",
	Run: func(_ *cobra.Command, _ []string) {
		ctx := context.Background()

		token := os.Getenv("GITEA_TOKEN")
		if token == "" {
			token = os.Getenv("GITHUB_TOKEN")
		}

		if token == "" {
			log.Fatal("GITEA_TOKEN or GITHUB_TOKEN environment variable must be set")
		}

		client := gitea.NewClient(*giteaURL, token, *owner, *repo)

		prs, err := client.ListOpenPRs(ctx)
		if err != nil {
			log.Fatalf("list PRs: %v", err)
		}

		if len(prs) == 0 {
			fmt.Println("No open pull requests.")

			return
		}

		for _, pr := range prs {
			fmt.Printf("#%-4d  %-60s  %s\n", pr.Number, pr.Title, pr.Head.Label)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkVersionCmd)
	rootCmd.AddCommand(GenerateYamlChartsCmd)
	rootCmd.AddCommand(createPRsCmd)
	rootCmd.AddCommand(listPRsCmd)
	rootCmd.AddCommand(mergePRCmd)

	rootCmd.PersistentFlags().String("versions", "versions.yaml", "versions.yaml file to user")

	debug = rootCmd.PersistentFlags().Bool("debug", false, "enable debug")
	filter = checkVersionCmd.Flags().String("filter", "", "filter to apply when checking helm/container images")

	giteaURL = rootCmd.PersistentFlags().String("gitea-url", "https://git.mkz.me", "Gitea instance base URL")
	owner = rootCmd.PersistentFlags().String("owner", "mycroft", "repository owner")
	repo = rootCmd.PersistentFlags().String("repo", "k8s-home", "repository name")
	baseBranch = createPRsCmd.Flags().String("base-branch", "main", "base branch for pull requests")
	dryRun = createPRsCmd.Flags().Bool("dry-run", false, "print what would be done without creating PRs")
	prFilter = createPRsCmd.Flags().String("filter", "", "filter to apply when checking helm/container images")
	prBranchOverride = createPRsCmd.Flags().String("branch", "", "override the generated branch name (requires exactly one update)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

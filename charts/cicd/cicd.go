package infra

import (
	"fmt"
	"os"

	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewCICDChart(scope constructs.Construct) cdk8s.Chart {
	chart := cdk8s.NewChart(
		scope,
		jsii.String("cicd"),
		&cdk8s.ChartProps{},
	)

	// builds kustomizations
	cicdYamlFile, err := k8s_helpers.BuildKustomize("./cicd")
	if err != nil {
		fmt.Printf("could not generate kustomization: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(cicdYamlFile)

	cdk8s.NewInclude(chart, jsii.String("cicd"), &cdk8s.IncludeProps{
		Url: jsii.String(cicdYamlFile),
	})

	return chart
}

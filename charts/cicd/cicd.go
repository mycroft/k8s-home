package infra

import (
	"fmt"
	"os"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewCICDChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "tekton-builds"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

	// builds kustomizations
	cicdYamlFile, err := kubehelpers.BuildKustomize("./cicd")
	if err != nil {
		fmt.Printf("could not generate kustomization: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(cicdYamlFile)

	// TODO: find a way to ensure all resources in charts are correctly namespaced.

	// namespace mentionned here is unused
	cdk8s.NewInclude(chart, jsii.String(namespace), &cdk8s.IncludeProps{
		Url: jsii.String(cicdYamlFile),
	})

	return chart
}

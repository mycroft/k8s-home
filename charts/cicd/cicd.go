package infra

import (
	"fmt"
	"os"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewCICDChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "tekton-builds"

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	// builds kustomizations
	cicdYamlFile, err := kubehelpers.BuildKustomize("./cicd")
	if err != nil {
		fmt.Printf("could not generate kustomization: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(cicdYamlFile)

	// TODO: find a way to ensure all resources in charts are correctly namespaced.

	// namespace mentionned here is unused
	cdk8s.NewInclude(chart.Cdk8sChart, jsii.String(namespace), &cdk8s.IncludeProps{
		Url: jsii.String(cicdYamlFile),
	})

	return chart
}

package infra

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewTektonChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	chart := builder.NewChart("tekton")

	versions := map[string]string{
		"pipelines": "v0.59.0",
		"triggers":  "v0.27.0",
		"dashboard": "v0.46.0",
	}

	cdk8s.NewInclude(chart.Cdk8sChart, jsii.String("tekton-pipeline"), &cdk8s.IncludeProps{
		Url: jsii.String(fmt.Sprintf("https://storage.googleapis.com/tekton-releases/pipeline/previous/%s/release.yaml", versions["pipelines"])),
	})

	cdk8s.NewInclude(chart.Cdk8sChart, jsii.String("tekton-triggers-release"), &cdk8s.IncludeProps{
		Url: jsii.String(fmt.Sprintf("https://storage.googleapis.com/tekton-releases/triggers/previous/%s/release.yaml", versions["triggers"])),
	})

	cdk8s.NewInclude(chart.Cdk8sChart, jsii.String("tekton-triggers-interceptors"), &cdk8s.IncludeProps{
		Url: jsii.String(fmt.Sprintf("https://storage.googleapis.com/tekton-releases/triggers/previous/%s/interceptors.yaml", versions["triggers"])),
	})

	cdk8s.NewInclude(chart.Cdk8sChart, jsii.String("tekton-dashboard"), &cdk8s.IncludeProps{
		Url: jsii.String(fmt.Sprintf("https://storage.googleapis.com/tekton-releases/dashboard/previous/%s/release.yaml", versions["dashboard"])),
	})

	return chart
}

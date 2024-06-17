package kubehelpers

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type Builder struct {
	App              cdk8s.App
	Context          context.Context
	Versions         Versions
	DockerImages     []string
	HelmRepositories map[string]string
}

type Chart struct {
	Cdk8sChart cdk8s.Chart
	Builder    *Builder
}

// NewBuilder creates a Builder context with cdk8s app, context & read versions file
func NewBuilder(ctx context.Context) *Builder {
	versions, err := ReadVersions()
	if err != nil {
		log.Fatal(err)
	}

	return &Builder{
		App:              cdk8s.NewApp(nil),
		Context:          ctx,
		Versions:         versions,
		HelmRepositories: make(map[string]string),
	}
}

// NewChart builds a cdk8s.Chart instance and returns it
func (builder *Builder) NewChart(namespace string) *Chart {
	return &Chart{
		Cdk8sChart: cdk8s.NewChart(
			builder.App,
			jsii.String(namespace),
			&cdk8s.ChartProps{},
		),
		Builder: builder,
	}
}

// BuildChart calls the passed callback with the current Builder context
func (builder *Builder) BuildChart(callback func(*Builder) *Chart) *Chart {
	return callback(builder)
}

// BuildChartLegacy calls the passed callback with the current Builder context (legacy version)
func (builder *Builder) BuildChartLegacy(callback func(context.Context, constructs.Construct) cdk8s.Chart) cdk8s.Chart {
	return callback(builder.Context, builder.App)
}

// RegisterContainerImage marks container images used for version checking
func (builder *Builder) RegisterContainerImage(image string) string {
	if val, exists := builder.Versions.Images[image]; exists {
		image = fmt.Sprintf("%s:%s", image, val)
	}

	builder.DockerImages = append(builder.DockerImages, image)

	return image
}

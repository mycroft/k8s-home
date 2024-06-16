package kubehelpers

import (
	"context"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type Builder struct {
	App      cdk8s.App
	Context  context.Context
	Versions Versions
}

type Chart struct {
	Cdk8sChart cdk8s.Chart
}

// NewBuilder creates a Builder context with cdk8s app, context & read versions file
func NewBuilder(ctx context.Context) *Builder {

	return &Builder{
		App:      cdk8s.NewApp(nil),
		Context:  ctx,
		Versions: versions,
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
	}
}

// BuildChart calls the passed callback with the current Builder context
func (builder *Builder) BuildChart(callback func(*Builder) cdk8s.Chart) Chart {
	return Chart{
		Cdk8sChart: callback(builder),
	}
}

// BuildChartLegacy calls the passed callback with the current Builder context (legacy version)
func (builder *Builder) BuildChartLegacy(callback func(context.Context, constructs.Construct) cdk8s.Chart) cdk8s.Chart {
	return callback(builder.Context, builder.App)
}

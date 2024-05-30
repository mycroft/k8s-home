package apps

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

const (
	whatismyipImage = "registry.mkz.me/mycroft/whatismyip:5"
)

func NewWhatIsMyIpChart(scope constructs.Construct) cdk8s.Chart {
	appName := "whatismyip"
	ingressHost := fmt.Sprintf("%s.services.mkz.me", appName)

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, appName)

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	k8s.NewKubeDeployment(
		chart,
		jsii.String("deploy"),
		&k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(appName),
			},
			Spec: &k8s.DeploymentSpec{
				Selector: &k8s.LabelSelector{
					MatchLabels: &labels,
				},
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels: &labels,
					},
					Spec: &k8s.PodSpec{
						Containers: &[]*k8s.Container{
							{
								Name:  jsii.String(appName),
								Image: jsii.String(whatismyipImage),
							},
						},
					},
				},
			},
		},
	)

	kubehelpers.NewAppIngress(
		chart,
		labels,
		appName,
		8080,
		ingressHost,
		map[string]string{},
	)

	return chart
}

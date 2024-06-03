package apps

import (
	"context"
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/imports/servicemonitor_monitoringcoreoscom"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

const (
	heyImage = "git.mkz.me/mycroft/hey:latest"
)

func NewHeyChart(ctx context.Context, scope constructs.Construct) cdk8s.Chart {
	namespace := "hey"
	appName := namespace
	appPort := 3000
	ingressHosts := []string{
		fmt.Sprintf("%s.services.mkz.me", appName),
		fmt.Sprintf("%s.mkz.cx", appName),
	}

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	kubehelpers.NewNamespace(chart, namespace)

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
								ImagePullPolicy: jsii.String("Always"),
								Name:            jsii.String(appName),
								Image:           jsii.String(heyImage),
							},
						},
					},
				},
			},
		},
	)

	kubehelpers.NewAppIngresses(
		chart,
		labels,
		appName,
		appPort,
		ingressHosts,
		"",
		map[string]string{},
	)

	servicemonitor_monitoringcoreoscom.NewServiceMonitor(
		chart,
		jsii.String("service-monitor"),
		&servicemonitor_monitoringcoreoscom.ServiceMonitorProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(namespace),
				Labels: &map[string]*string{
					"release": jsii.String("prometheus"),
				},
			},
			Spec: &servicemonitor_monitoringcoreoscom.ServiceMonitorSpec{
				Selector: &servicemonitor_monitoringcoreoscom.ServiceMonitorSpecSelector{
					MatchLabels: &map[string]*string{},
				},
				Endpoints: &[]*servicemonitor_monitoringcoreoscom.ServiceMonitorSpecEndpoints{
					{
						Path: jsii.String("/metrics"),
						Port: jsii.String("http"),
					},
				},
			},
		},
	)

	return chart
}

package kubehelpers

import (
	"git.mkz.me/mycroft/k8s-home/imports/servicemonitor_monitoringcoreoscom"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type ServiceMonitor struct {
	Labels          map[string]string
	ServicePortName string
}

func (chart *Chart) NewServiceMonitor(serviceMonitor *ServiceMonitor) {
	if chart.Namespace == "" {
		panic("namespace was not defined")
	}

	servicePortName := serviceMonitor.ServicePortName
	if servicePortName == "" {
		servicePortName = "http"
	}

	labels := ToLabelsPtr(serviceMonitor.Labels)

	servicemonitor_monitoringcoreoscom.NewServiceMonitor(
		chart.Cdk8sChart,
		jsii.String("service-monitor"),
		&servicemonitor_monitoringcoreoscom.ServiceMonitorProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(chart.Namespace),
				Labels: &map[string]*string{
					"release": jsii.String("prometheus"),
				},
			},
			Spec: &servicemonitor_monitoringcoreoscom.ServiceMonitorSpec{
				Selector: &servicemonitor_monitoringcoreoscom.ServiceMonitorSpecSelector{
					MatchLabels: &labels,
				},
				Endpoints: &[]*servicemonitor_monitoringcoreoscom.ServiceMonitorSpecEndpoints{
					{
						Path: jsii.String("/metrics"),
						Port: jsii.String(servicePortName),
					},
				},
			},
		},
	)
}

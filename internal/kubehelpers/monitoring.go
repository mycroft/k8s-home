package kubehelpers

import (
	"git.mkz.me/mycroft/k8s-home/imports/scrapeconfig_monitoringcoreoscom"
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

// ScrapeTargetBearerToken references the key of a Secret, in the kube-prometheus-stack
// chart's namespace, holding the bearer token used to authenticate scrape requests.
type ScrapeTargetBearerToken struct {
	SecretName string
	SecretKey  string
}

// ScrapeTarget describes an external HTTP endpoint that Prometheus should scrape
// directly, for hosts that aren't backed by a Kubernetes Service/Endpoints object
// (e.g. bare-metal exporters reachable only by hostname).
type ScrapeTarget struct {
	// Name is used as both the ScrapeConfig resource name and the "job" label.
	Name string
	// Target is the host:port to scrape, e.g. "moonstone.lan.mkz.me:9100".
	Target string
	// Path is the HTTP path to scrape metrics from. Defaults to "/metrics".
	Path string
	// Scheme is the HTTP scheme used to scrape the target. Defaults to HTTP.
	Scheme scrapeconfig_monitoringcoreoscom.ScrapeConfigSpecScheme
	// BearerToken optionally authenticates the scrape request. Leave nil for
	// unauthenticated targets.
	BearerToken *ScrapeTargetBearerToken
}

// CreateScrapeConfig registers a prometheus-operator ScrapeConfig so kube-prometheus-stack
// scrapes an external static target, bypassing Kubernetes Service/Endpoints discovery.
// The resource is created in chart's namespace with the `release: prometheus` label
// required by the chart's scrapeConfigSelector.
func CreateScrapeConfig(chart *Chart, target ScrapeTarget) {
	path := target.Path
	if path == "" {
		path = "/metrics"
	}

	scheme := target.Scheme
	if scheme == "" {
		scheme = scrapeconfig_monitoringcoreoscom.ScrapeConfigSpecScheme_HTTP
	}

	var authorization *scrapeconfig_monitoringcoreoscom.ScrapeConfigSpecAuthorization
	if target.BearerToken != nil {
		authorization = &scrapeconfig_monitoringcoreoscom.ScrapeConfigSpecAuthorization{
			Credentials: &scrapeconfig_monitoringcoreoscom.ScrapeConfigSpecAuthorizationCredentials{
				Name: jsii.String(target.BearerToken.SecretName),
				Key:  jsii.String(target.BearerToken.SecretKey),
			},
		}
	}

	scrapeconfig_monitoringcoreoscom.NewScrapeConfig(
		chart.Cdk8sChart,
		jsii.String(target.Name),
		&scrapeconfig_monitoringcoreoscom.ScrapeConfigProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(chart.Namespace),
				Labels: &map[string]*string{
					"release": jsii.String("prometheus"),
				},
			},
			Spec: &scrapeconfig_monitoringcoreoscom.ScrapeConfigSpec{
				StaticConfigs: &[]*scrapeconfig_monitoringcoreoscom.ScrapeConfigSpecStaticConfigs{
					{
						Targets: &[]*string{
							jsii.String(target.Target),
						},
						Labels: &map[string]*string{
							"job": jsii.String(target.Name),
						},
					},
				},
				MetricsPath:   jsii.String(path),
				Scheme:        scheme,
				Authorization: authorization,
			},
		},
	)
}

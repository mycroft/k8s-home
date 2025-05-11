package kubehelpers

import (
	"context"
	"fmt"
	"strings"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type AppServiceOption struct {
	Name string
}

// NewAppService creates a Service resource for the application.
func NewAppService(
	chart cdk8s.Chart,
	namespace string,
	serviceName string,
	labels map[string]*string,
	portName string,
	portNumber uint,
	opts ...AppServiceOption,
) k8s.KubeService {
	metadata := k8s.ObjectMeta{
		Namespace: jsii.String(namespace),
	}

	spec := k8s.ServiceSpec{
		Ports: &[]*k8s.ServicePort{
			{
				Name: jsii.String(portName),
				Port: jsii.Number(float64(portNumber)),
			},
		},
		Selector: &labels,
	}

	for _, opt := range opts {
		if opt.Name != "" {
			metadata.Name = jsii.String(opt.Name)
		}
	}

	props := k8s.KubeServiceProps{
		Metadata: &metadata,
		Spec:     &spec,
	}

	return k8s.NewKubeService(
		chart,
		jsii.String(serviceName),
		&props,
	)
}

type AppIngressOption struct {
	Name string
}

// NewAppIngresses creates one or more Ingress resources for the application.
func NewAppIngresses(
	ctx context.Context,
	chart cdk8s.Chart,
	labels map[string]*string,
	namespace string,
	appPort uint,
	ingressHosts []string,
	serviceName string,
	customAnnotations map[string]string,
	opts ...AppIngressOption,
) {
	annotations := map[string]*string{
		"cert-manager.io/cluster-issuer": jsii.String("letsencrypt-prod"),
	}

	for k, v := range customAnnotations {
		annotations[k] = jsii.String(v)
	}

	portName := "http"

	if serviceName == "" {
		if ContextGetDebug(ctx) {
			fmt.Printf("creating a service for %s' Ingress (%s)\n", *chart.ToString(), strings.Join(ingressHosts, ", "))
		}

		svc := NewAppService(
			chart,
			namespace,
			"svc",
			labels,
			portName,
			appPort,
		)

		serviceName = *svc.Name()
	}

	rules := []*k8s.IngressRule{}
	hosts := []*string{}

	for _, ingressHost := range ingressHosts {
		rules = append(rules, &k8s.IngressRule{
			Host: jsii.String(ingressHost),
			Http: &k8s.HttpIngressRuleValue{
				Paths: &[]*k8s.HttpIngressPath{
					{
						Backend: &k8s.IngressBackend{
							Service: &k8s.IngressServiceBackend{
								Name: jsii.String(serviceName),
								Port: &k8s.ServiceBackendPort{
									Name: jsii.String(portName),
								},
							},
						},
						Path:     jsii.String("/"),
						PathType: jsii.String("Prefix"),
					},
				},
			},
		})

		hosts = append(hosts, jsii.String(ingressHost))
	}

	metadatas := k8s.ObjectMeta{
		Annotations: &annotations,
		Namespace:   jsii.String(namespace),
	}

	for _, opt := range opts {
		if opt.Name != "" {
			metadatas.Name = jsii.String(opt.Name)
		}
	}

	k8s.NewKubeIngress(
		chart,
		jsii.String("ingress"),
		&k8s.KubeIngressProps{
			Metadata: &metadatas,
			Spec: &k8s.IngressSpec{
				IngressClassName: jsii.String("traefik"),
				Rules:            &rules,
				Tls: &[]*k8s.IngressTls{
					{
						Hosts:      &hosts,
						SecretName: jsii.String("secret-tls-www"),
					},
				},
			},
		},
	)
}

// NewAppIngress creates an Ingress resource for the application.
func NewAppIngress(
	ctx context.Context,
	chart cdk8s.Chart,
	labels map[string]*string,
	namespace string,
	appPort uint,
	ingressHost string,
	serviceName string,
	customAnnotations map[string]string,
) {
	NewAppIngresses(
		ctx,
		chart,
		labels,
		namespace,
		appPort,
		[]string{
			ingressHost,
		},
		serviceName,
		customAnnotations,
	)
}

type Ingress struct {
	Namespace   string
	Name        string
	Port        uint
	Ingresses   []string
	Labels      map[string]*string
	ServiceName string
	Annotations map[string]string
}

// NewIngress creates an ingress attached to deployment name/port
// If ingress.ServiceName is unset, a new service will be created
func (chart *Chart) NewIngress(ingress *Ingress) {
	if ingress.ServiceName == "" {
		ingress.ServiceName = ingress.Name

		NewAppService(
			chart.Cdk8sChart,
			ingress.Namespace,
			ingress.ServiceName,
			ingress.Labels,
			"http",
			ingress.Port,
			AppServiceOption{
				Name: ingress.Name,
			},
		)
	}

	NewAppIngresses(
		chart.Builder.Context,
		chart.Cdk8sChart,
		ingress.Labels,
		ingress.Namespace,
		ingress.Port,
		ingress.Ingresses,
		ingress.ServiceName,
		ingress.Annotations,
		AppIngressOption{
			Name: ingress.Name,
		},
	)
}

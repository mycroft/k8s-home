package kubehelpers_test

import (
	"testing"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// ingressPaths returns the paths routed by the chart's Ingresses.
func ingressPaths(t *testing.T, chart *kubehelpers.Chart) []string {
	t.Helper()

	paths := []string{}

	for _, object := range *chart.Cdk8sChart.ToJson() {
		manifest, _ := object.(map[string]interface{})
		if manifest["kind"] != "Ingress" {
			continue
		}

		spec, _ := manifest["spec"].(map[string]interface{})
		rules, _ := spec["rules"].([]interface{})

		for _, rule := range rules {
			http, _ := rule.(map[string]interface{})["http"].(map[string]interface{})
			entries, _ := http["paths"].([]interface{})

			for _, entry := range entries {
				path, _ := entry.(map[string]interface{})["path"].(string)
				paths = append(paths, path)
			}
		}
	}

	return paths
}

// Not parallel: concurrent first calls into the jsii runtime race while it
// starts its kernel process, and crash.
func TestNewIngressPath(t *testing.T) {
	cases := []struct {
		name string
		path string
		want string
	}{
		{"default routes the whole host", "", "/"},
		{"a path routes only that prefix", "/webhooks", "/webhooks"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			builder := &kubehelpers.Builder{
				App:              cdk8s.NewApp(nil),
				Versions:         kubehelpers.Versions{},
				HelmRepositories: map[string]string{},
			}
			chart := builder.NewChart("app")
			chart.NewNamespace("app")

			chart.NewIngress(&kubehelpers.Ingress{
				Name:        "app",
				Ingresses:   []string{"app.services.mkz.me"},
				ServiceName: "app",
				Path:        tc.path,
			})

			got := ingressPaths(t, chart)
			if len(got) != 1 || got[0] != tc.want {
				t.Errorf("ingress paths = %v, want [%s]", got, tc.want)
			}
		})
	}
}

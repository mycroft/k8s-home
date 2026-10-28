package kubehelpers_test

import (
	"reflect"
	"testing"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// Not parallel: concurrent first calls into the jsii runtime race while it
// starts its kernel process, and crash.
func TestNewAppDeploymentSecurityContext(t *testing.T) {
	cases := []struct {
		name string
		opt  kubehelpers.AppDeploymentOption
		// nil when the pod must have no securityContext
		want map[string]interface{}
	}{
		{"no option sets none", kubehelpers.AppDeploymentOption{}, nil},
		{"RunAsNonRoot for a numeric image user", kubehelpers.AppDeploymentOption{RunAsNonRoot: true}, map[string]interface{}{"runAsNonRoot": true}},
		{"RunAsUser for a named image user", kubehelpers.AppDeploymentOption{RunAsNonRoot: true, RunAsUser: 1000}, map[string]interface{}{"runAsNonRoot": true, "runAsUser": float64(1000)}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			chart := cdk8s.NewChart(cdk8s.NewApp(nil), jsii.String("app"), nil)
			labels := map[string]*string{"app.kubernetes.io/name": jsii.String("app")}

			kubehelpers.NewAppDeployment(chart, "app", "app", "org/app:1.0.0", labels, nil, nil, nil, tc.opt)

			manifests := *chart.ToJson()
			if len(manifests) != 1 {
				t.Fatalf("got %d manifests, want 1 Deployment", len(manifests))
			}

			deployment, _ := manifests[0].(map[string]interface{})
			spec, _ := deployment["spec"].(map[string]interface{})
			template, _ := spec["template"].(map[string]interface{})
			podSpec, _ := template["spec"].(map[string]interface{})
			got, _ := podSpec["securityContext"].(map[string]interface{})

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("spec.template.spec.securityContext = %v, want %v", got, tc.want)
			}
		})
	}
}

// Not parallel, for the same reason as above.
func TestNewAppDeploymentPodAnnotations(t *testing.T) {
	chart := cdk8s.NewChart(cdk8s.NewApp(nil), jsii.String("app"), nil)
	labels := map[string]*string{"app.kubernetes.io/name": jsii.String("app")}

	kubehelpers.NewAppDeployment(chart, "app", "app", "org/app:1.0.0", labels, nil, nil, nil, kubehelpers.AppDeploymentOption{
		PodAnnotations: map[string]string{"configMapHash": "abc"},
	})

	deployment, _ := (*chart.ToJson())[0].(map[string]interface{})
	spec, _ := deployment["spec"].(map[string]interface{})
	template, _ := spec["template"].(map[string]interface{})
	metadata, _ := template["metadata"].(map[string]interface{})
	got, _ := metadata["annotations"].(map[string]interface{})

	want := map[string]interface{}{"configMapHash": "abc"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("spec.template.metadata.annotations = %v, want %v", got, want)
	}
}

package kubehelpers_test

import (
	"testing"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// Not parallel: concurrent first calls into the jsii runtime race while it
// starts its kernel process, and crash.
func TestCreateHelmReleaseCRDsPolicy(t *testing.T) {
	cases := []struct {
		name        string
		opts        []kubehelpers.HelmReleaseOption
		wantInstall interface{}
		wantUpgrade interface{}
	}{
		// Flux's defaults: create missing CRDs on install, skip them on upgrade.
		{"default leaves the policy to Flux", nil, nil, nil},
		{"WithReplaceCRDs replaces on install and upgrade", []kubehelpers.HelmReleaseOption{kubehelpers.WithReplaceCRDs()}, "CreateReplace", "CreateReplace"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			builder := &kubehelpers.Builder{
				App:              cdk8s.NewApp(nil),
				Versions:         kubehelpers.Versions{HelmCharts: map[string]string{"repo/chart": "1.0.0"}},
				HelmRepositories: map[string]string{},
			}

			release := builder.NewChart("app").CreateHelmRelease("app", "repo", "chart", "app", tc.opts...)

			spec, _ := release.ToJson().(map[string]interface{})["spec"].(map[string]interface{})
			install, _ := spec["install"].(map[string]interface{})
			upgrade, _ := spec["upgrade"].(map[string]interface{})

			if got := install["crds"]; got != tc.wantInstall {
				t.Errorf("spec.install.crds = %v, want %v", got, tc.wantInstall)
			}

			if got := upgrade["crds"]; got != tc.wantUpgrade {
				t.Errorf("spec.upgrade.crds = %v, want %v", got, tc.wantUpgrade)
			}

			if got := install["skipCRDs"]; got != false {
				t.Errorf("spec.install.skipCRDs = %v, want false", got)
			}
		})
	}
}

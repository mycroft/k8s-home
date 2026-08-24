package kubehelpers_test

import (
	"testing"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

func TestLatestMatchingVersion(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		candidates []string
		pattern    string
		want       string
	}{
		{"picks highest semver", []string{"1.2.3", "2.0.0", "1.9.0"}, ".+", "2.0.0"},
		{"restores v prefix", []string{"v1.2.3", "v2.0.0"}, ".+", "v2.0.0"},
		{"pattern filters non-semver", []string{"1.2.3", "1.2.3-alpine"}, `^[0-9]+\.[0-9]+(\.[0-9]+)?$`, "1.2.3"},
		{"no match returns empty", []string{"1.2.3-alpine"}, `^[0-9]+\.[0-9]+$`, ""},
		{"empty candidates", nil, ".+", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := kubehelpers.LatestMatchingVersion(tc.candidates, tc.pattern); got != tc.want {
				t.Errorf("LatestMatchingVersion() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestOciPrereleaseTags(t *testing.T) {
	got := kubehelpers.OciPrereleaseTags([]string{"1.2.3", "2.0.0-rc.1", "1.10.0", "not-a-version", "2.0.0-alpha"})

	want := []string{"1.2.3", "1.10.0"}

	if len(got) != len(want) {
		t.Fatalf("OciPrereleaseTags() = %v, want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("OciPrereleaseTags()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

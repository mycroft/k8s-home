package kubehelpers_test

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
)

// ExampleOciPrereleaseTags shows how OCI chart tags are reduced to stable
// releases before version comparison.
func ExampleOciPrereleaseTags() {
	tags := []string{"1.2.3", "2.0.0-rc.1", "1.10.0", "not-a-version"}

	stable := kubehelpers.OciPrereleaseTags(tags)

	fmt.Println(stable)

	// Output: [1.2.3 1.10.0]
}

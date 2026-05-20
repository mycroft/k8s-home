package kubehelpers

import "github.com/aws/jsii-runtime-go"

func ToLabelsPtr(labels map[string]string) map[string]*string {
	newLabels := map[string]*string{}

	for k, v := range labels {
		newLabels[k] = jsii.String(v)
	}

	return newLabels
}

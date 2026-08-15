package apps

import "git.mkz.me/mycroft/k8s-home/internal/kubehelpers"

func NewLlamaChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	return builder.NewExternalApplicationChart(kubehelpers.ExternalApplicationConfig{
		Name:          "llama",
		Hostname:      "llm-api.iop.cx",
		ExternalName:  "10.0.0.7",
		Port:          80,
		CertificateID: "openjunk-cert",
	})
}

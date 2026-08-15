package apps

import "git.mkz.me/mycroft/k8s-home/internal/kubehelpers"

func NewMusicAssistantChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	return builder.NewExternalApplicationChart(kubehelpers.ExternalApplicationConfig{
		Name:         "music-assistant",
		Hostname:     "music-assistant.iop.cx",
		ExternalName: "10.0.0.7",
		Port:         8095,
	})
}

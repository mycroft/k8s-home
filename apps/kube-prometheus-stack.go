package apps

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
)

func NewKubePrometheusStackChart(scope constructs.Construct) cdk8s.Chart {
	appName := "kube-prometheus-stack"
	namespace := "monitoring"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(appName),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"prometheus-community",
		"https://prometheus-community.github.io/helm-charts",
	)

	k8s_helpers.NewSealedSecret(
		chart,
		namespace,
		"grafana-secret",
		map[string]*string{
			"password": jsii.String("AgA8Gj167Qfm8VlG1VkeU86rMF2LIotZmXjTHBw+t3bTY3zUjjIDJaJQKqYySTekDpGgA/+vR7erScMbVreyD+5rDYDe5XxOvrpK4Wa3C0QZ1GXhT9jg9T9DopQ7EwJQ7HUwmIOop3PIz7wa7SACabdDUioaxxhpI/gAAYWdA7jl25fXOm7AvdtxjmHaxzxY6AjrZ4VhJqVm8rcsHVF0J3Xw83HqproIM7JBg/oqavRrmLev6nbwUEE15sKXstfGqrXkhOeb/Gw8e8Ic+DzC1OuhfZChgohyLC/FtLx7TAymtb85DpXkxW3MaRU0bvna53dfA1tchTsi4SK/UMP/xQZgZemStTntoycEZ6XvKbLPjMdjuBza6ThZOfCEVSMBjqKKuaWNh3VnjN5s//vJxHw/WG4jeHwfntgo8dp1us69IQ/HV76tgKUABBBPoTGRtaN+jUSE4KTONcvIqtzc/JgM1EAxkYWoWXvLyNg3UOjPD4Jjr6xxxFvc8DC8HK85Xy91s8Yd1umcjYZ08VrkfSY3OVoS54S+dCy365cu5bZ4klAphs24DBj7Kmd16VQduqcIPwUSdRAUjuQg2qyPg+W78lysKZE/Gs8H203C4x1yE0K2tN1zVW4KQeGH276q+y9RWYDVQw9Op7vzZGMNkYjBH947EPt6uD/UHx5hdLsjLDvY5OMZhQaZ35dEyD9RxqM1ymjeUHEloUVp1dHNJ2nU"),
			"user":     jsii.String("AgBYdmNCWw4ai42G50Fl1ydHNzKMTWfcpJGKqNvTaNUxe79EWQY0/fvUi5nWWzLNTOJpNHVgISPIdcW+WZT4Cu7ZPvmMzynEcqDMjaFWkOlbXtCrOQTfIOKEX1vAeGB6NWlLflpzqGykgiaUUXZ7xI7qvHv7ooUfLi/e7LTgZxDfqBC6DFrJ3R0YprS+PUoqHZTfkr15I4dnz2foA1A50XpLMweSqrmqcR1eHUSduqubplr5xIfUtdOhnuhCCe3H6PyjEoicCZ8cBPqbhHCAg2TwDifZrn9efJXkqT0WheZn17p2jPuUeCA+RrtOe8NNXAFm0+1zeowzVJCVCCone0b23SJ1Syre1muiTWEitD1Cs78xAn/2iAUWDUqgolrgV6LgeU8zwwt2013tbzd1ykzJv09tEjwzG2l49bv9+3L2lX18IwQgUGZToJhB9MQt/3z9h0sO6D2WelqLNpCA7GEywSCMIZLhszR9pTUIVPUZRxCQzgun51iySeGWiC1XLfF8WSPcG3YrmLWq5SWZI+xA4vRRHLtOyTSOubYAjF6gx2R/Rd4pMgIU4V/YnqrnlWDpnZz6oNQ6dHDM+vsRSRMQdJ7IJ4hdxvKMPzflvy1GptbqCijG1f/3bV13OWGhyWX/biHCIhcH4rxIQLbfuqs+BeFuWFdOUoSvB2d4HRM4wlEPoQNuCMhbtesfKFxN/uDpcGvhZQ=="),
		},
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"prometheus-community",  // repoName; must be in flux-system
		"kube-prometheus-stack", // chart name
		"prometheus",            // release name
		"41.3.1",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"kube-prometheus-stack.yaml",
			),
		},
		nil,
	)

	return chart
}

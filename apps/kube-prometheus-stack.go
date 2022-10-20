package apps

import (
	"os"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
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

	contents, err := os.ReadFile("configs/kube-prometheus-stack.yaml")
	if err != nil {
		panic(err)
	}

	cm := k8s.NewKubeConfigMap(
		chart,
		jsii.String("cm"),
		&k8s.KubeConfigMapProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
			},
			Data: &map[string]*string{
				"kube-prometheus-stack.yaml": jsii.String(string(contents)),
			},
		},
	)

	k8s_helpers.NewSealedSecret(
		chart,
		namespace,
		"grafana-secret",
		map[string]*string{
			"password": jsii.String("AgCCADVAPJC4DoSw0A/tI1sjYMRm0KC68tx2exEA+NgNj63vlWjj1mgbZjBMDRil7z/knMQqwuMv39G10I/RpJa7Sady/YmnCWkaWQmfdW6cMeemFxGqy1ARGlbT5l1G4cYQW8bViBC32g9CAJhXVzYxsIcO1qZgkr43aJeXzWvT2VhN7whDD6T4gCg1/JolW7jHu7blGdUP/qLeSirFn/TevxTFbpU/ebvy9ixUuqaIbDD1dg40k/1C8HNnFA3I7AV6VuOVE8GbdxWmS8Rr86okBw4F0u2FQebmIx2wxoDfyeJIY7qZix6a8yYqC6t5c+1JAQ41UKdNlYqlGCBZ4OPqNaJHLO9JcLhTPbv4SpeFS0vtWriF5wMQXEqlRH6Bo7vP0rdKpic08PVkSRk2yK5013XTDNI2lZU1/Ck4UKmTnIV7e7bjgiWsMfgv1EbRmLeEHMhuOhEjS5z0gUVjWacxXBhZEgz/nHypjplVP2fOlXNAt/QMOrz17d7Y19PrITRutp2gs0tpIkepXV9Oz96frpz+VSfDxbbsumX+BrTKqXpDiHj46UagOIn1H0BTgpuFO/cnOnzIcyRjoVoNm0Anc+J3n4VCuN6dMYjmRQfip3aNNf+8o1Be9K2lzxrWsyJ7t0LKO6UTFIzsTmCTNRiOg17WsrSUA5K1dVho02cRnIcZI+IYmC26Uy72xXMU4OLGJjPVzFaWUOiQptvxEBWZ"),
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
			{
				Name:    *cm.Name(),
				KeyName: "kube-prometheus-stack.yaml",
			},
		},
		map[string]*string{
			"configMapHash": jsii.String(k8s_helpers.ComputeConfigMapHash(cm)),
		},
	)

	return chart
}

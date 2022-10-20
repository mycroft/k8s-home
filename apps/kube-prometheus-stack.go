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
			"password": jsii.String("AgANbKKJq+SiBRbBWnjizT1F3pBLm5teM9w/G+/YAw8QuQfvH/5/hPmQozxxG5DHP1AUJ12yJMw6MDvipXd0fVELEs1fMzPUlL2sPksYvwvlHYOTCXWOCrjIQ6q4NsivEg5wWGJs1KgvbkehmJvhCqgp47/x5BgCYtp0J92FNfpS4o51x74Ae3Pza/mAAH+EeIE0x1MkMeAz+o73o8DBdgxRU7Vhnv3UdCVCWgZQBNuufpEN1dMWcEQ1yUs0yMDnOB3uiMqtTx+7WQWUu/66FMzPGsnEer/lk2IYR+1DlVCL8AWsl1ED3TIBDI7L8UpzArpzPjANkKxLWQg5mJLXa+elb2dSNK/sQN/zEOKVp8VpIiP8XIz/8HXAyS5A6ustcWu/ieDHhtAJNDvvWdzER04LqPfOhFW89C1SM43IATeN8burBM60O+1safD0v1aPZUhd3te9Etfs0J+J4NnPhvwH0OGspM+SuWtRnMuh5VwEFhK2Q/JZM0JRkulFaswu6DPco4N3ftDZYdrabBntnJnHvdxCcqWWHU8xIQ0Jn0U+T5ZHhmpGYih6TOKiv5AXzpF9HSsJBKBwwsFDCZ8QnDJ8jFF++CSEMIh8nqlJwjkG3rMKhqDapanDh3alzW07CUUl6B4P7MBjjYrUfLohrCi9dvBAWj5pP04BuUu4u/r/xdruurdSNaZRmk5msoBQOlFqx09eJ+8iNk0WWfHbp9J1"),
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

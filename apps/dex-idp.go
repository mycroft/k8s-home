package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewDexIdpChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "dex-idp"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"dex",
		"https://charts.dexidp.io",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"dex", // repo name
		"dex", // chart name
		"dex", // release name
		"0.12.1",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"dex-idp.yaml",
			),
		},
		nil,
	)

	// Create configuration
	// See contrib/extract-dex-config.sh to extract/seal new configuration
	// This is used in configs/dex-idp.yaml
	k8s_helpers.NewSealedSecret(
		chart,
		namespace,
		"dex-config",
		map[string]*string{
			"config.yaml": jsii.String("AgDIMrmuMJFfWPxdtE89v7ISUKbb3M11VsfqWgEPjN3HdlTg82QXskAA+EVt1D0uAunlsdImaU+8xF2dvO+Ei4sED34HvIQJe9AD6dcNp8uVrsMM3TUa2K8tS8LVT9XcGjG2J0iKMTz8fQsL/3j9WFWb7KCIwpJZxEJfI4GyhWVVbtLi4/7ImplCwW1OUNxJOn78R6Gfng3e+Q2NSiNcBnwt6kGgclLrsHoJJ5n0wE16TEeWeZBewRep+fwhRD/ijN0MYtsqNZDhue6wVODwWCee63yhpsj+UMv7/+XPqxxMLvHRITaAepqZNBs6gZGl9BWNWaEGAjfkf2EqTNnG6ubH4SZOgS910IV5MjkkLIQx5jZvKWn1eEWCLe8ECNGav7fe9m6eruyLoYoKmRO7eW6B/i0g2nK9o0yVwlMLilEILllLwHXfUEfHyN59/HTmfzVk9K/jXQlSYZnl+i6CIwW+yDyhROuJUpA8BvYATX8vQbQseu1Jw3A11uPkdHQWQvupVUQkLDIzeZzgDkYHSYQ37XepEaI38IHg6hh9i8XJ9tGPoGYhsy9YMNBLBZKfEoogwIZU3hnoKT8YHMB07fNs3ku3SGmmiENJsF3HAsKU9EvoO9EJW9B02A5Iux0qZu8nl2IUSfbXmzdXflQ4qBd1gPhIQ2WNxc5ceLCP1or1KZtmdqmxSJun2n30KYitoZW/s4+vWyuEYpM1QUw4XvlxvGnzSBKpoenkGJx9yOHRPWoZD+xjpaPt2eS9c4EeKJ8QbIsWrPZs3VWFU13Cp+y+kJwUEwQ14TLKUHx+sUIehfJheRXTxW9WEUB2HlsY9ANHay5rPbk="),
		},
	)

	return chart
}

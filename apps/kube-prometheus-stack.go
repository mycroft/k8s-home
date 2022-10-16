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
			"password": jsii.String("AgCwmVbMGtTjnwKpO52JW7JFhhnyds07aymDrqXb6FZqqzDaRWKHRWH7h2grXhsgOtfcLlA36+4GugxhUUjoJElO/eOI3gM41t6Hk+XjuFBbRpswKvC7dGXpS0ofRZq51bwjZXaCHo9AjOXTovz4Jg7rIKGwg4gY4CV+hcRhlmdGLtKaUglRfguc1K1iRvNAM5yf406fCSFPU3SJm6REBl+3npN+2hBD/CiYx+PYz+Qy297WOtAosXagA7NJi/qtt+lpURuvmNkrCmCg+E/JHysg6lRegbK6wz6rAA+a1z2GE3RcZno1Qx7cfovfVqmEpkRkzyQJrbub+T9iblJaIiqYriEKV98wFq5om+8N5rZBkme5b1U+WaoZyfH9NfNEsas9/tdYlvCE8ESNcuLjSLXwE/fLplUF3DbbLrX+QIe+cyIdcp708VgxUoQAnndvDSCObyWAGyqwaHRhdfuL+UsGQIttT1wyjEmlejoJxE/ttHtQspiAaVB/CTonfuM840jAuQrtnZasDStC5kgfaNwyL8pbV143FWFcQZHTl8K38jSfhOYYC2Z9iLRqieaCZMdNt5hAi7F+BkdYETLnf5FNqRTqdlM5yj63KbDCPH3AuyfO15+lxZ9/sF9z0SQWQZTNMukN0GgF1keJ3mMgZJm1qI6A4FydWKiAIer+1t6zEHYkEFZ/dVpSSMhX6QpCB+LMxULVXbldWcE="),
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

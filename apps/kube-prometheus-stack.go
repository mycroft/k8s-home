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
			"passwordKey": jsii.String("AgC7dhPNyPCSH/TWRSe66PN+yQI03dHersP3HOVsoJMCr+WPLrhU2L92wZ07Q5ibWSTTdx//lxFptA6L+kkLHOIWXD3gmGx63sFQcz88zC4NXfaxyc6XA2CdFyinMl/qjmdS76QLa7/Yfw1tTfseMhcC5Mcy/dpxH5iZlembJK6vh3YQZ8j2ppfwnfRsMvU9V0q2UN1bxEl/JNzEpaxsD9+fGEfQciCY/5Kg201C5JkiFOBVdvC8486cmhLSgrG4BT5Lfhekdg7f4DGySE+aZMVyxFML3qhWPFw2xvrwB8cATNaNaBgliDp6NsU3wynYj44yZd/Zbe0jD/pilg/AhoFndnmLad/X3tS3TVnkGAOzgDtQwkBM4hSLvjPzY/yYHVl/uDaw/B1WedqQz3Qqb0XX17Q90T71Sthyq68S74w+0rXHb3g1P1s2d9tK21CjJrIYxBnxtwZrfQzOHgIO5gKq//HGN2edKwOCcZcAzH08eEy5l9F1lD0haujT02AEyUIBHhsOaLhEzg058SlZ5aZASpVZbr8HLn/BBsKhpTpv8FAYNA5CJCL+IbbtHfWywNidreu/ifxWn/Q4Ir+lyuZ022yhMWNTULBWONK5OptuO2ETnTWW2CwbS5lCSGEaOSMj9uaONmfGxbzR7ayxldygvvhC7YwrOnHbyOQU6/xCiaFf6d5a5BgvOhaUt96I9WZLPWRw0SN/IvhdcuKdqB8cWHSHfpgj1CO1"),
			"userKey":     jsii.String("AgBOl1LNXPIEtspMvfFIxuAy4YYUWG+HyEfyJ/169BtVY7Jcgt1hQsib0tH8VYJE/qEc1hOUr+Q892pf6BBlRS4IbbczSRNP//YmTss3QK6xg53AJENsHvtVKpmnfOZIAf030M7HPbCLUwmv5SWH1OMnCPG4JVqVEMn9/2q4mBQbafQ9yzgveNxIzJWz5dexzmuG6tKHa7P6aYmRbeLiI2517ARxXlSJhunoETW20k/sYUueSQ1Ai6WHKP7JbPJ8W7rKWI57cck8IM41/X3RzPuf7sVKVCYBeA5q63vh2psWYP7ffLx9w4Yucv8xCksp4d/XD+Pr7GRIL1X0cchr7GfQx58WSeoTL5sQLPZCXMKyzfl+51n197YiB48ElgiWrDyxmJ3H60g7MDePT85AYvmulUXjhUD80bIN/c5xENypNRm61JWTmfDO0XkyEBBbFN6gd3KVZ1KDgS9kCQEB5px/+/d3kHajRQmulF5bNS6O4SSy3gY0YjK9US3cYVp3YiCcmROMN/PQQZnFogXZjstgZVqNba6yJSoNQPMEOoIQLLPExJ2NNtXz2HADNgQthvxQ14hdE51CXFU4Sb8EhFDSlKrns1JWei+mc09mPDEztmTsO4Xmj6CSxm3MVhFzHUqsxWD1MUcK1vGncMtoa+6nNZH27BrWEYBXEBErZ3qL6kWBM0xZQoPKGqpMO+ynf/FQJRjvKg=="),
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

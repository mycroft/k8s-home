package infra

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewCapacitorChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "flux-system"
	repositoryName := "onechart"

	chart := cdk8s.NewChart(
		scope,
		jsii.String("capacitor"),
		&cdk8s.ChartProps{},
	)

	kubehelpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://chart.onechart.dev",
	)

	kubehelpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName,
		"onechart",
		"capacitor",
		map[string]string{},
		[]kubehelpers.HelmReleaseConfigMap{
			kubehelpers.CreateHelmValuesConfig(
				chart,
				namespace,
				repositoryName,
				"capacitor.yaml",
			),
		},
		nil,
	)

	k8s.NewKubeServiceAccount(
		chart,
		jsii.String("sa"),
		&k8s.KubeServiceAccountProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("capacitor"),
			},
		},
	)

	k8s.NewKubeClusterRole(
		chart,
		jsii.String("clusterrole"),
		&k8s.KubeClusterRoleProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
				Name:      jsii.String("capacitor"),
			},
			Rules: &[]*k8s.PolicyRule{
				{
					ApiGroups: jsii.Strings(
						"networking.k8s.io",
						"apps",
						"",
					),
					Resources: jsii.Strings(
						"pods",
						"pods/log",
						"ingresses",
						"deployments",
						"services",
						"secrets",
						"events",
						"configmaps",
					),
					Verbs: jsii.Strings(
						"get",
						"watch",
						"list",
					),
				},
				{
					ApiGroups: jsii.Strings(
						"source.toolkit.fluxcd.io",
						"kustomize.toolkit.fluxcd.io",
						"helm.toolkit.fluxcd.io",
					),
					Resources: jsii.Strings(
						"gitrepositories",
						"ocirepositories",
						"buckets",
						"helmrepositories",
						"helmcharts",
						"kustomizations",
						"helmreleases",
					),
					Verbs: jsii.Strings(
						"get",
						"watch",
						"list",
						"patch",
					),
				},
			},
		},
	)

	k8s.NewKubeClusterRoleBinding(
		chart,
		jsii.String("clusterrolebinding"),
		&k8s.KubeClusterRoleBindingProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String("capacitor"),
			},
			Subjects: &[]*k8s.Subject{
				{
					Kind:      jsii.String("ServiceAccount"),
					Name:      jsii.String("capacitor"),
					Namespace: jsii.String(namespace),
				},
			},
			RoleRef: &k8s.RoleRef{
				Kind:     jsii.String("ClusterRole"),
				Name:     jsii.String("capacitor"),
				ApiGroup: jsii.String("rbac.authorization.k8s.io"),
			},
		},
	)

	return chart
}

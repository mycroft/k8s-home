package apps

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewKubernetesDashboardChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "kubernetes-dashboard"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		"kubernetes-dashboard",
		"https://kubernetes.github.io/dashboard/",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		"kubernetes-dashboard",
		"kubernetes-dashboard", // chart name
		"kubernetes-dashboard", // release name
		"5.11.0",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				"kubernetes-dashboard.yaml",
			),
		},
		nil,
	)

	// Create a Service Account & ClusterRoleBinding
	// https://github.com/kubernetes/dashboard/blob/master/docs/user/access-control/creating-sample-user.md
	k8s.NewKubeServiceAccount(
		chart,
		jsii.String("sa"),
		&k8s.KubeServiceAccountProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("admin"),
				Namespace: jsii.String(namespace),
			},
		},
	)

	k8s.NewKubeClusterRoleBinding(
		chart,
		jsii.String("cluster-role-binding-admin"),
		&k8s.KubeClusterRoleBindingProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String("admin"),
			},
			RoleRef: &k8s.RoleRef{
				ApiGroup: jsii.String("rbac.authorization.k8s.io"),
				Kind:     jsii.String("ClusterRole"),
				Name:     jsii.String("cluster-admin"),
			},
			Subjects: &[]*k8s.Subject{
				{
					Kind:      jsii.String("ServiceAccount"),
					Name:      jsii.String("admin"),
					Namespace: jsii.String(namespace),
				},
			},
		},
	)
	return chart
}

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
	repositoryName := "kubernetes-dashboard"
	chartName := "kubernetes-dashboard"
	releaseName := "kubernetes-dashboard"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)

	k8s_helpers.CreateHelmRepository(
		chart,
		repositoryName,
		"https://kubernetes.github.io/dashboard/",
	)

	k8s_helpers.CreateHelmRelease(
		chart,
		namespace,
		repositoryName,
		chartName,
		releaseName,
		"6.0.6",
		map[string]string{},
		[]k8s_helpers.HelmReleaseConfigMap{
			k8s_helpers.CreateHelmValuesConfig(
				chart,
				namespace,
				repositoryName,
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

	// https://kubernetes.io/docs/concepts/configuration/secret/#service-account-token-secrets
	// The token will be automatically filled in this secret:
	k8s.NewKubeSecret(
		chart,
		jsii.String("sa-secret"),
		&k8s.KubeSecretProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String("secret-sa"),
				Namespace: jsii.String(namespace),
				Annotations: &map[string]*string{
					"kubernetes.io/service-account.name": jsii.String("admin"),
				},
			},
			Type: jsii.String("kubernetes.io/service-account-token"),
		},
	)

	k8s.NewKubeClusterRole(
		chart,
		jsii.String("cluster-role-read-only"),
		&k8s.KubeClusterRoleProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String("cluster-role-read-only"),
			},
			Rules: &[]*k8s.PolicyRule{
				{
					ApiGroups: &[]*string{
						jsii.String("*"),
					},
					Verbs: &[]*string{
						jsii.String("get"),
						jsii.String("list"),
						jsii.String("watch"),
					},
					Resources: &[]*string{
						jsii.String("namespaces"),
						jsii.String("events"),
						jsii.String("pods"),
						jsii.String("cronjobs"),
						jsii.String("jobs"),
						jsii.String("replicasets"),
						jsii.String("deployments"),
						jsii.String("daemonsets"),
						jsii.String("replicationcontrollers"),
						jsii.String("statefulsets"),
						jsii.String("ingresses"),
						jsii.String("ingressesclasses"),
						jsii.String("services"),
						jsii.String("configmaps"),
						jsii.String("persistentvolumeclaims"),
						jsii.String("storageclasses"),
						jsii.String("clusterrolebindings"),
						jsii.String("clusterroles"),
						jsii.String("networkpolicies"),
						jsii.String("nodes"),
						jsii.String("persistentvolumes"),
						jsii.String("rolebindings"),
						jsii.String("roles"),
						jsii.String("serviceaccounts"),
						jsii.String("customresourcedefinitions"),
					},
				},
			},
		},
	)

	k8s.NewKubeClusterRoleBinding(
		chart,
		jsii.String("cluster-role-binding-kubernetes-dashboard"),
		&k8s.KubeClusterRoleBindingProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String("kubernetes-dashboard"),
			},
			RoleRef: &k8s.RoleRef{
				ApiGroup: jsii.String("rbac.authorization.k8s.io"),
				Kind:     jsii.String("ClusterRole"),
				Name:     jsii.String("cluster-role-read-only"),
			},
			Subjects: &[]*k8s.Subject{
				{
					Kind:      jsii.String("ServiceAccount"),
					Name:      jsii.String("kubernetes-dashboard"),
					Namespace: jsii.String(namespace),
				},
			},
		},
	)

	return chart
}

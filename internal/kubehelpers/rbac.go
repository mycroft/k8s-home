package kubehelpers

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/jsii-runtime-go"
)

// NewReadAllButSecretsClusterRole grants serviceAccount, in the chart's
// namespace, get/list/watch on everything in the cluster but Secrets, through
// a ClusterRole and its binding, both named read-all-but-secrets-<namespace>.
// It is meant for the in-cluster MCP servers agent-nexus reads the cluster
// through.
//
// RBAC cannot subtract, so the rules list what is granted:
//   - the core group by name, leaving out secrets and the subresources that
//     reach into pods and nodes (exec, attach, portforward, proxy), which a
//     "*" would grant with get;
//   - the other groups with "*": none of them has such a subresource. Leaves
//     out generators.external-secrets.io, whose generatorstates hold
//     generated credentials.
//
// A newly installed API group stays unreadable until added here; kubectl
// api-resources --verbs=list -o name lists them.
func (chart *Chart) NewReadAllButSecretsClusterRole(serviceAccount string) {
	if chart.Namespace == "" {
		panic("namespace was not defined")
	}

	name := "read-all-but-secrets-" + chart.Namespace
	readVerbs := jsii.Strings("get", "list", "watch")

	k8s.NewKubeClusterRole(
		chart.Cdk8sChart,
		jsii.String("read-all-but-secrets"),
		&k8s.KubeClusterRoleProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String(name),
			},
			Rules: &[]*k8s.PolicyRule{
				{
					ApiGroups: jsii.Strings(""),
					Resources: jsii.Strings(
						"componentstatuses",
						"configmaps",
						"endpoints",
						"events",
						"limitranges",
						"namespaces",
						"nodes",
						"persistentvolumeclaims",
						"persistentvolumes",
						"pods",
						"pods/log",
						"podtemplates",
						"replicationcontrollers",
						"resourcequotas",
						"serviceaccounts",
						"services",
					),
					Verbs: readVerbs,
				},
				{
					ApiGroups: jsii.Strings(
						"acme.cert-manager.io",
						"admissionregistration.k8s.io",
						"apiextensions.k8s.io",
						"apiregistration.k8s.io",
						"apps",
						"aquasecurity.github.io",
						"autoscaling",
						"batch",
						"bitnami.com",
						"certificates.k8s.io",
						"cert-manager.io",
						"configuration.konghq.com",
						"coordination.k8s.io",
						"deuxfleurs.fr",
						"discovery.k8s.io",
						"events.k8s.io",
						"external-secrets.io",
						"flowcontrol.apiserver.k8s.io",
						"gateway.networking.k8s.io",
						"helm.cattle.io",
						"helm.toolkit.fluxcd.io",
						"hub.traefik.io",
						"infra.mkz.me",
						"k3s.cattle.io",
						"k8s.mkz.me",
						"kustomize.toolkit.fluxcd.io",
						"kyverno.io",
						"longhorn.io",
						"metrics.k8s.io",
						"mkz.me",
						"monitoring.coreos.com",
						"monitoring.grafana.com",
						"networking.k8s.io",
						"node.k8s.io",
						"notification.toolkit.fluxcd.io",
						"policy",
						"postgresql.cnpg.io",
						"rbac.authorization.k8s.io",
						"resource.k8s.io",
						"rollout-operator.grafana.com",
						"scheduling.k8s.io",
						"scylla.scylladb.com",
						"snapshot.storage.k8s.io",
						"source.toolkit.fluxcd.io",
						"storage.k8s.io",
						"traefik.io",
						"velero.io",
						"wgpolicyk8s.io",
					),
					Resources: jsii.Strings("*"),
					Verbs:     readVerbs,
				},
			},
		},
	)

	k8s.NewKubeClusterRoleBinding(
		chart.Cdk8sChart,
		jsii.String("read-all-but-secrets-binding"),
		&k8s.KubeClusterRoleBindingProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String(name),
			},
			Subjects: &[]*k8s.Subject{
				{
					Kind:      jsii.String("ServiceAccount"),
					Name:      jsii.String(serviceAccount),
					Namespace: jsii.String(chart.Namespace),
				},
			},
			RoleRef: &k8s.RoleRef{
				Kind:     jsii.String("ClusterRole"),
				Name:     jsii.String(name),
				ApiGroup: jsii.String("rbac.authorization.k8s.io"),
			},
		},
	)
}

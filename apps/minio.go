package apps

import (
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewMinio(scope constructs.Construct) cdk8s.Chart {
	namespace := "minio"

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)

	k8s_helpers.CreateExternalSecret(chart, namespace, "minio-tenant")

	// XXX TODO: Create the tenant here
	/*
			apiVersion: minio.min.io/v2
		kind: Tenant
		metadata:
		  name: minio-storage
		  namespace: minio
		spec:
		  credsSecret:
		    name: minio-tenant
		  pools:
		    - servers: 4
		      name: pool
		      volumesPerServer: 2
		      volumeClaimTemplate:
		        metadata:
		          name: data
		        spec:
		          storageClassName: longhorn-crypto-global
		          accessModes:
		            - ReadWriteOnce
		          resources:
		            requests:
		              storage: 2Gi
		      containerSecurityContext:
		        runAsUser: 1000
		        runAsGroup: 1000
		        runAsNonRoot: true
	*/
	// XXX Configure TLS certs?

	return chart
}

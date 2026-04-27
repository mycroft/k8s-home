package apps

// Note that the "outline" user in postgres needs to CREATE EXTENSION. Therefore the user was set up with the following permission:
// ALTER USER outline WITH SUPERUSER;

import (
	"fmt"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewOutlineChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	appName := "outline"
	namespace := appName
	appIngress := fmt.Sprintf("%s.iop.cx", appName)

	outlineImage := builder.RegisterContainerImage("outlinewiki/outline")

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "postgres")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "secret")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "oidc")
	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "garage")

	_, redisServiceName := chart.NewRedisStatefulset(namespace)

	env := []*k8s.EnvVar{
		{
			Name: jsii.String("DATABASE_URL"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("postgres"),
					Key:  jsii.String("url"),
				},
			},
		},
		{
			Name:  jsii.String("REDIS_URL"),
			Value: jsii.String(fmt.Sprintf("redis://%s:6379", redisServiceName)),
		},
		{
			Name: jsii.String("SECRET_KEY"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("secret"),
					Key:  jsii.String("secret_key"),
				},
			},
		},
		{
			Name: jsii.String("UTILS_SECRET"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("secret"),
					Key:  jsii.String("utils_secret"),
				},
			},
		},
		{
			Name:  jsii.String("URL"),
			Value: jsii.String(fmt.Sprintf("https://%s", appIngress)),
		},
		{
			Name:  jsii.String("COLLABORATION_URL"),
			Value: jsii.String(fmt.Sprintf("https://%s", appIngress)),
		},
		{
			Name: jsii.String("OIDC_CLIENT_ID"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("oidc"),
					Key:  jsii.String("client_id"),
				},
			},
		},
		{
			Name: jsii.String("OIDC_CLIENT_SECRET"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("oidc"),
					Key:  jsii.String("client_secret"),
				},
			},
		},
		{
			Name: jsii.String("OIDC_AUTH_URI"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("oidc"),
					Key:  jsii.String("auth_uri"),
				},
			},
		},
		{
			Name: jsii.String("OIDC_TOKEN_URI"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("oidc"),
					Key:  jsii.String("token_uri"),
				},
			},
		},
		{
			Name: jsii.String("OIDC_USERINFO_URI"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("oidc"),
					Key:  jsii.String("userinfo_uri"),
				},
			},
		},
		{
			Name: jsii.String("OIDC_LOGOUT_URI"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("oidc"),
					Key:  jsii.String("logout_uri"),
				},
			},
		},
		{
			Name: jsii.String("AWS_ACCESS_KEY_ID"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("garage"),
					Key:  jsii.String("access_key_id"),
				},
			},
		},
		{
			Name: jsii.String("AWS_SECRET_ACCESS_KEY"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("garage"),
					Key:  jsii.String("secret_access_key"),
				},
			},
		},
		{
			Name:  jsii.String("AWS_REGION"),
			Value: jsii.String("garage"),
		},
		{
			Name:  jsii.String("AWS_S3_UPLOAD_BUCKET_URL"),
			Value: jsii.String("https://s3-api.iop.cx"),
		},
		{
			Name:  jsii.String("AWS_S3_UPLOAD_BUCKET_NAME"),
			Value: jsii.String("outline"),
		},
		{
			Name:  jsii.String("AWS_S3_FORCE_PATH_STYLE"),
			Value: jsii.String("true"),
		},
		{
			Name:  jsii.String("FILE_STORAGE_UPLOAD_MAX_SIZE"),
			Value: jsii.String("268435456"),
		},
	}

	outlineLabels := map[string]*string{
		"app.kubernetes.io/component": jsii.String(appName),
	}

	appPort := uint(3000)

	_, svcName := kubehelpers.NewStatefulSet(chart.Cdk8sChart, kubehelpers.StatefulSetConfig{
		Namespace: namespace,
		AppName:   appName,
		AppImage:  outlineImage,
		AppPort:   appPort,
		Labels:    outlineLabels,
		Env:       env,
		Storages: []kubehelpers.StatefulSetVolume{
			{
				Name:        "data",
				MountPath:   "/var/lib/outline/data",
				StorageSize: "16Gi",
			},
		},
	})

	kubehelpers.NewAppIngress(
		builder.Context,
		chart.Cdk8sChart,
		outlineLabels,
		appName,
		appPort,
		appIngress,
		svcName,
		map[string]string{},
	)

	return chart
}

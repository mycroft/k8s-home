package apps

import (
	"log"
	"os"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"git.mkz.me/mycroft/k8s-home/internal/kubehelpers"
	"github.com/aws/jsii-runtime-go"
)

func NewPrivatebinChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
	namespace := "privatebin"
	appName := namespace
	appPort := 8080
	appIngress := "privatebin.services.mkz.me"
	privatebinImage := kubehelpers.RegisterDockerImage("privatebin/nginx-fpm-alpine")

	chart := builder.NewChart(namespace)
	chart.NewNamespace(namespace)

	kubehelpers.CreateSecretStore(chart.Cdk8sChart, namespace)

	kubehelpers.CreateExternalSecret(chart.Cdk8sChart, namespace, "minio")

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	confPHP, err := os.ReadFile("configs/privatebin/conf.php")
	if err != nil {
		log.Fatalf("Could not read privatebin configuration file: %s", err.Error())
	}

	configMap := k8s.NewKubeConfigMap(
		chart.Cdk8sChart,
		jsii.String("config"),
		&k8s.KubeConfigMapProps{
			Metadata: &k8s.ObjectMeta{
				Namespace: jsii.String(namespace),
			},
			Data: &map[string]*string{
				"conf.php": jsii.String(string(confPHP)),
			},
		},
	)

	env := []*k8s.EnvVar{
		{
			Name: jsii.String("AWS_ACCESS_KEY_ID"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("minio"),
					Key:  jsii.String("access_key_id"),
				},
			},
		},
		{
			Name: jsii.String("AWS_SECRET_ACCESS_KEY"),
			ValueFrom: &k8s.EnvVarSource{
				SecretKeyRef: &k8s.SecretKeySelector{
					Name: jsii.String("minio"),
					Key:  jsii.String("secret_access_key"),
				},
			},
		},
		{
			Name:  jsii.String("CONFIG_PATH"),
			Value: jsii.String("/run/cfg"),
		},
	}

	cmd := []string{
		"mkdir /run/cfg",
		"cp /srv/cfg/conf.php /run/cfg/conf.php",
		"sed -i \"s/AWS_ACCESS_KEY_ID/$AWS_ACCESS_KEY_ID/\" /run/cfg/conf.php",
		"sed -i \"s/AWS_SECRET_ACCESS_KEY/$AWS_SECRET_ACCESS_KEY/\" /run/cfg/conf.php",
		"/etc/init.d/rc.local",
	}

	kubehelpers.NewAppDeployment(
		chart.Cdk8sChart,
		namespace,
		appName,
		privatebinImage,
		labels,
		env,
		cmd,
		[]kubehelpers.ConfigMapMount{
			{
				Name:      "config",
				ConfigMap: configMap,
				MountPath: "/srv/cfg",
			},
		},
	)

	kubehelpers.NewAppIngress(
		builder.Context,
		chart.Cdk8sChart,
		labels,
		appName,
		appPort,
		appIngress,
		"",
		map[string]string{},
	)

	return chart
}

package apps

import (
	"io/ioutil"
	"log"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	k8s_helpers "git.mkz.me/mycroft/k8s-home/k8s-helpers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewPrivatebinChart(scope constructs.Construct) cdk8s.Chart {
	namespace := "privatebin"
	appName := namespace
	appPort := 8080
	appIngress := "privatebin.services.mkz.me"
	privatebinImage := k8s_helpers.RegisterDockerImage("privatebin/nginx-fpm-alpine")

	chart := cdk8s.NewChart(
		scope,
		jsii.String(namespace),
		&cdk8s.ChartProps{},
	)

	k8s_helpers.NewNamespace(chart, namespace)
	k8s_helpers.CreateSecretStore(chart, namespace)

	k8s_helpers.CreateExternalSecret(chart, namespace, "minio")

	labels := map[string]*string{
		"app.kubernetes.io/name": jsii.String(appName),
	}

	confPHP, err := ioutil.ReadFile("configs/privatebin/conf.php")
	if err != nil {
		log.Fatalf("Could not read privatebin configuration file: %s", err.Error())
	}

	configMap := k8s.NewKubeConfigMap(
		chart,
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

	k8s_helpers.NewAppDeployment(
		chart,
		namespace,
		appName,
		privatebinImage,
		labels,
		env,
		cmd,
		[]k8s_helpers.ConfigMapMount{
			{
				Name:      "config",
				ConfigMap: configMap,
				MountPath: "/srv/cfg",
			},
		},
	)

	k8s_helpers.NewAppIngress(
		chart,
		labels,
		appName,
		appPort,
		appIngress,
		map[string]string{},
	)

	return chart
}

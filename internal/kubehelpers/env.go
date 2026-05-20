package kubehelpers

import (
	"git.mkz.me/mycroft/k8s-home/imports/k8s"
	"github.com/aws/jsii-runtime-go"
)

type EnvValueFromSecret struct {
	Key  string
	Name string
}

type EnvValue struct {
	Value           string
	ValueFromSecret EnvValueFromSecret
}

type EnvEntry struct {
	Name  string
	Value EnvValue
}

func ToK8sEnv(env []EnvEntry) []*k8s.EnvVar {
	returnedEnv := []*k8s.EnvVar{}

	for _, envItem := range env {
		if envItem.Value.Value != "" {
			returnedEnv = append(returnedEnv, &k8s.EnvVar{
				Name:  jsii.String(envItem.Name),
				Value: jsii.String(envItem.Value.Value),
			})
		}

		if envItem.Value.ValueFromSecret.Key != "" {
			returnedEnv = append(returnedEnv, &k8s.EnvVar{
				Name: jsii.String(envItem.Name),
				ValueFrom: &k8s.EnvVarSource{
					SecretKeyRef: &k8s.SecretKeySelector{
						Key:  jsii.String(envItem.Value.ValueFromSecret.Key),
						Name: jsii.String(envItem.Value.ValueFromSecret.Name),
					},
				},
			})
		}
	}

	return returnedEnv
}

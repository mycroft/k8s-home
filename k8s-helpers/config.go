package k8s_helpers

import (
	"crypto/sha256"
	"fmt"
	"sort"

	"git.mkz.me/mycroft/k8s-home/imports/k8s"
)

// ComputeConfigMapHash will build an unique hash according to the content
// of the given configMap files/keys. It can be re-used like the following
// in PodTemplate in Deployments objects:
//
//	PodMetadata: &cdk8s.ApiObjectMetadata{
//	  Annotations: &map[string]*string{
//	    "configMapHash": jsii.String(avxlib.ComputeConfigMapHash(configMap)),
//	  }
//	}
//
// This can be re-used to enforce redeploying pods automatically if their given
// configMap changes but not the deployment as it (except the given annotation).
func ComputeConfigMapHash(configMaps ...k8s.KubeConfigMap) string {
	h := sha256.New()

	for _, cm := range configMaps {
		jsonContent := cm.ToJson().(map[string]interface{})

		keys := make([]string, 0)
		for k := range jsonContent["data"].(map[string]interface{}) {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		dataContent := jsonContent["data"].(map[string]interface{})

		for _, k := range keys {
			h.Write([]byte(dataContent[k].(string)))
		}

	}

	/*
		                }



		                for _, k := range keys {
		                        h.Write([]byte(*values[k]))
		                }
		        }

		        return fmt.Sprintf("%x", h.Sum(nil))

				json := cm.ToJson().(map[string]interface{})

				for k, v := range json["data"].(map[string]interface{}) {
					fmt.Printf("%v / %v", k, v.(string))
				}
	*/

	return fmt.Sprintf("%x", h.Sum(nil))
}

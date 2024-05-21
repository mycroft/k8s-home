package k8s_helpers

import (
	"fmt"
	"os"

	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// Build using kustomize a kustomization path and returns file path to built yaml file.
// On error, file will be deleted and an empty string along the error will be returned.
//
// do not forget to cleanup with defer os.Remove(returnedFilePath)
func BuildKustomize(kustomizationDirPath string) (string, error) {
	fSys := filesys.MakeFsOnDisk()

	f, err := os.CreateTemp("/tmp", "kustomizaion-output-*.yaml")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary yaml file: %v", err)

	}
	defer f.Close()

	yamlFilePath := f.Name()

	k := krusty.MakeKustomizer(&krusty.Options{
		PluginConfig: &types.PluginConfig{},
	})
	m, err := k.Run(fSys, kustomizationDirPath)
	if err != nil {
		defer os.Remove(yamlFilePath)
		return "", fmt.Errorf("failed to run kustomize: %v", err)
	}

	yml, err := m.AsYaml()
	if err != nil {
		defer os.Remove(yamlFilePath)
		return "", fmt.Errorf("failed to generate yaml: %v", err)
	}

	if err = fSys.WriteFile(yamlFilePath, yml); err != nil {
		defer os.Remove(yamlFilePath)
		return "", fmt.Errorf("failed to write yaml in temporary file: %v", err)
	}

	// filePath can be included in a cdk8s chart by:
	//
	// cdk8s.NewInclude(chart, jsii.String("my-chart"), &cdk8s.IncludeProps{
	// 	Url: jsii.String(filePath),
	// })

	return yamlFilePath, nil
}

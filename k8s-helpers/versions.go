package k8s_helpers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/Masterminds/semver"
	"gopkg.in/yaml.v2"
)

func GetRepoIndex(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

type Index struct {
	ApiVersion string              `yaml:"apiVersion"`
	Entries    map[string][]*Entry `yaml:"entries"`
}

type Entry struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Annotations map[string]string `yaml:"annotations"`
}

func GetEntriesFromIndex(body []byte) (map[string][]*Entry, error) {
	var index Index
	err := yaml.Unmarshal(body, &index)
	if err != nil {
		return map[string][]*Entry{}, err
	}

	return index.Entries, nil
}

// CheckVersions checks the installed releases for update
func CheckVersions() {
	for _, helmRelease := range helmChartVersions {
		if _, ok := helmRepositories[helmRelease.RepositoryName]; !ok {
			panic(fmt.Sprintf("Unknown repo %s", helmRelease.RepositoryName))
		}

		repositoryUrl := helmRepositories[helmRelease.RepositoryName]

		body, err := GetRepoIndex(fmt.Sprintf("%s/index.yaml", repositoryUrl))
		if err != nil {
			panic(err)
		}

		entries, err := GetEntriesFromIndex(body)
		if err != nil {
			panic(err)
		}

		// find entries for this chart
		if _, ok := entries[helmRelease.ChartName]; !ok {
			panic(fmt.Sprintf("No chart for name %s", helmRelease.ChartName))
		}

		versions := []string{}
		for _, chartVersion := range entries[helmRelease.ChartName] {
			if _, ok := chartVersion.Annotations["artifacthub.io/prerelease"]; ok {
				if chartVersion.Annotations["artifacthub.io/prerelease"] == "true" {
					continue
				}
			}
			versions = append(versions, chartVersion.Version)
		}

		hasVPrefix := false

		semvers := make([]*semver.Version, len(versions))
		for i, version := range versions {
			if version[0] == 'v' {
				version = version[1:]
				hasVPrefix = true
			}
			v, err := semver.NewVersion(version)
			if err != nil {
				panic(fmt.Sprintf("Invalid version %s: %s", version, err))
			}
			semvers[i] = v
		}

		sort.Sort(sort.Reverse(semver.Collection(semvers)))

		lastVersion := ""
		if hasVPrefix {
			lastVersion = "v" + semvers[0].Original()
		} else {
			lastVersion = semvers[0].Original()
		}

		if lastVersion != helmRelease.Version {
			fmt.Printf("%s: %s => %s\n", helmRelease.ChartName, helmRelease.Version, lastVersion)
		}
	}
}

package k8s_helpers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
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

	fmt.Printf("Helm charts updates:\n")

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

	fmt.Printf("Docker images:\n")

	for _, image := range dockerImages {
		parts := strings.Split(image, ":")
		GetLastImageTag(parts[0], parts[1])
	}
}

var dockerImages = []string{}

func RegisterDockerImage(image string) string {
	dockerImages = append(dockerImages, image)

	return image
}

func GetLastImageTag(image, version string) {
	// Set the image name and registry URL
	imageName := image
	// registryURL := "https://registry-1.docker.io"

	// Create a new registry client
	ref, err := name.ParseReference(fmt.Sprintf("%s:%s", imageName, "latest"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	tags, err := remote.ListWithContext(ctx, ref.Context())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: %s\n", image, version)

	if version == "latest" || version == "" {
		return
	}

	v, err := semver.NewVersion(version)
	if err != nil {
		panic(fmt.Sprintf("Invalid version %s: %s", version, err))
	}

	for _, tag := range tags {
		v2, err := semver.NewVersion(tag)
		if err != nil {
			continue
		}

		if v2.GreaterThan(v) {
			fmt.Printf("New version: %s\n", v2)
		}
	}
}

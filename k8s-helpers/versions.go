package k8s_helpers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

func GetHelmUpdates() (map[string]string, error) {
	retVersions := map[string]string{}

	for _, helmRelease := range helmChartVersions {
		chartName := fmt.Sprintf("%s/%s", helmRelease.RepositoryName, helmRelease.ChartName)
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
			retVersions[chartName] = fmt.Sprintf("%s;%s", helmRelease.Version, lastVersion)
		}
	}

	return retVersions, nil
}

// CheckVersions checks the installed releases for update
func CheckVersions() {
	helmVersions, err := GetHelmUpdates()
	if err != nil {
		panic(err)
	}

	for k, v := range helmVersions {
		fmt.Printf("%s;%s\n", k, v)
	}

	for _, image := range dockerImages {
		parts := strings.Split(image, ":")

		versions := GetLastImageTag(parts[0], parts[1])

		if len(versions) > 0 {
			fmt.Printf("%s;%s;%s\n", parts[0], parts[1], versions[len(versions)-1])
		}

	}
}

var dockerImages = []string{}

func RegisterDockerImage(image string) string {
	ReadVersions()

	if val, exists := versions.Images[image]; exists {
		image = fmt.Sprintf("%s:%s", image, val)
	}

	dockerImages = append(dockerImages, image)

	return image
}

func GetLastImageTag(image, version string) []string {
	retVersions := []string{}
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

	if version == "latest" || version == "" {
		return retVersions
	}

	v, err := semver.NewVersion(version)
	if err != nil {
		panic(fmt.Sprintf("Invalid version %s: %s", version, err))
	}

	for _, tag := range tags {
		if strings.HasPrefix(image, "linuxserver") && !strings.HasPrefix(tag, "v") {
			continue
		}

		v2, err := semver.NewVersion(tag)
		if err != nil {
			continue
		}

		if v2.GreaterThan(v) {
			if strings.HasPrefix(image, "linuxserver") {
				retVersions = append(retVersions, tag)
			} else {
				retVersions = append(retVersions, v2.String())
			}
		}
	}

	return retVersions
}

type Versions struct {
	HelmCharts map[string]string `yaml:"helmcharts"`
	Images     map[string]string `yaml:"images"`
}

var versions = Versions{}

func ReadVersions() {
	if len(versions.Images) != 0 || len(versions.HelmCharts) != 0 {
		return
	}

	body, err := os.ReadFile("versions.yaml")
	if err != nil {
		panic(fmt.Sprintf("Could not open versions.yaml: %s", err.Error()))
	}

	err = yaml.Unmarshal(body, &versions)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal versions.yaml: %s", err.Error()))
	}
}

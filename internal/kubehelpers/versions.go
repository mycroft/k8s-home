package kubehelpers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"gopkg.in/yaml.v2"
)

type Index struct {
	APIVersion string              `yaml:"apiVersion"`
	Entries    map[string][]*Entry `yaml:"entries"`
}

type Entry struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Annotations map[string]string `yaml:"annotations"`
}

// Patterns is part of versions.yaml configuration file structure
type Patterns struct {
	HelmCharts map[string]string `yaml:"helmcharts"`
	Images     map[string]string `yaml:"images"`
}

// Versions is part of versions.yaml configuration file structure
type Versions struct {
	HelmCharts map[string]string `yaml:"helmcharts"`
	Images     map[string]string `yaml:"images"`
	Patterns   Patterns          `yaml:"patterns"`
}

func GetRepoIndex(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetEntriesFromIndex(body []byte) (map[string][]*Entry, error) {
	var index Index
	err := yaml.Unmarshal(body, &index)
	if err != nil {
		return map[string][]*Entry{}, err
	}

	return index.Entries, nil
}

func (builder *Builder) GetHelmUpdates(debug bool, filter string) (map[string]string, error) {
	retVersions := map[string]string{}
	versions := builder.Versions

	for _, helmRelease := range helmChartVersions {
		chartName := fmt.Sprintf("%s/%s", helmRelease.RepositoryName, helmRelease.ChartName)
		if _, ok := builder.HelmRepositories[helmRelease.RepositoryName]; !ok {
			panic(fmt.Sprintf("Unknown repo %s", helmRelease.RepositoryName))
		}

		if filter != "" && !strings.Contains(chartName, filter) {
			if debug {
				log.Printf("skipping helm chart check: %s", chartName)
			}
			continue
		}

		repositoryURL := builder.HelmRepositories[helmRelease.RepositoryName]

		if strings.HasPrefix(repositoryURL, "oci://") {
			if debug {
				log.Printf("skipped %s as oci:// is not supported", repositoryURL)
			}
			continue
		}

		body, err := GetRepoIndex(fmt.Sprintf("%s/index.yaml", repositoryURL))
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

		foundVersions := []string{}
		for _, chartVersion := range entries[helmRelease.ChartName] {
			if _, ok := chartVersion.Annotations["artifacthub.io/prerelease"]; ok {
				if chartVersion.Annotations["artifacthub.io/prerelease"] == "true" {
					continue
				}
			}

			pattern := ".+"
			if _, ok := versions.Patterns.HelmCharts[chartName]; ok {
				pattern = versions.Patterns.HelmCharts[chartName]
			}

			matched, err := regexp.MatchString(pattern, chartVersion.Version)
			if err != nil {
				panic(err)
			}

			if !matched {
				continue
			}

			foundVersions = append(foundVersions, chartVersion.Version)
		}

		hasVPrefix := false

		semvers := make([]*semver.Version, len(foundVersions))
		for i, version := range foundVersions {
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
func (builder *Builder) CheckVersions(debug bool, filter string) {
	searched := make(map[string]bool)
	versions := builder.Versions

	helmVersions, err := builder.GetHelmUpdates(debug, filter)
	if err != nil {
		panic(err)
	}

	for k, v := range helmVersions {
		fmt.Printf("%s;%s\n", k, v)
	}

	for _, image := range builder.DockerImages {
		parts := strings.Split(image, ":")

		if _, ok := searched[parts[0]]; ok {
			continue
		}

		searched[parts[0]] = true

		if filter != "" && !strings.Contains(image, filter) {
			if debug {
				log.Printf("skipping image check: %s...", image)
			}
			continue
		}

		if debug {
			log.Printf("checking %s...", image)
		}

		startTs := time.Now()

		pattern := ".+"
		if _, ok := versions.Patterns.Images[parts[0]]; ok {
			pattern = versions.Patterns.Images[parts[0]]
		}

		foundVersions := GetLastImageTag(debug, parts[0], parts[1], pattern)

		if len(foundVersions) > 0 {
			fmt.Printf("%s;%s;%s\n", parts[0], parts[1], foundVersions[len(foundVersions)-1])
		}

		doneTs := time.Now()

		if debug {
			log.Printf("done checking after %d msec", doneTs.UnixMilli()-startTs.UnixMilli())
		}

	}
}

func GetLastImageTag(debug bool, image, version, pattern string) []string {
	retVersions := []string{}
	imageName := image

	// Create a new registry client
	ref, err := name.ParseReference(fmt.Sprintf("%s:%s", imageName, "latest"))
	if err != nil {
		panic(err)
	}

	tags, err := remote.List(ref.Context(), remote.WithContext(context.Background()))
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

		matched, err := regexp.MatchString(pattern, tag)
		if err != nil {
			panic(err)
		}

		if !matched {
			if debug {
				log.Printf("tag %s was skipped as it does not match pattern %s", tag, pattern)
			}
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

// ReadVersions reads versions.yaml, parses it and return its configuration in a Version instance
func ReadVersions(versionsFile string) (Versions, error) {
	var versions Versions

	if len(versions.Images) != 0 || len(versions.HelmCharts) != 0 {
		return Versions{}, nil
	}

	body, err := os.ReadFile(versionsFile)
	if err != nil {
		return Versions{}, fmt.Errorf("could not open versions.yaml: %s", err.Error())
	}

	err = yaml.Unmarshal(body, &versions)
	if err != nil {
		return Versions{}, fmt.Errorf("could not unmarshal versions.yaml: %s", err.Error())
	}

	return versions, nil
}

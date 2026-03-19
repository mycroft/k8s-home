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
	"sync"
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
	versions := builder.Versions

	// Collect unique repo URLs to fetch, validating that all repos are known.
	repoURLs := make(map[string]string) // repoName -> URL

	for _, helmRelease := range helmChartVersions {
		if _, ok := builder.HelmRepositories[helmRelease.RepositoryName]; !ok {
			panic(fmt.Sprintf("Unknown repo %s", helmRelease.RepositoryName))
		}

		repoURL := builder.HelmRepositories[helmRelease.RepositoryName]
		if !strings.HasPrefix(repoURL, "oci://") {
			repoURLs[helmRelease.RepositoryName] = repoURL
		}
	}

	// Fetch all repo indexes concurrently.
	var mu sync.Mutex

	repoEntries := make(map[string]map[string][]*Entry) // repoURL -> entries

	var wg sync.WaitGroup

	for _, repoURL := range repoURLs {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			body, err := GetRepoIndex(fmt.Sprintf("%s/index.yaml", url))
			if err != nil {
				panic(err)
			}

			entries, err := GetEntriesFromIndex(body)
			if err != nil {
				panic(err)
			}

			mu.Lock()
			repoEntries[url] = entries
			mu.Unlock()
		}(repoURL)
	}

	wg.Wait()

	// Process each chart against the cached indexes.
	retVersions := make(map[string]string)

	for _, helmRelease := range helmChartVersions {
		chartName := fmt.Sprintf("%s/%s", helmRelease.RepositoryName, helmRelease.ChartName)

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

		entries, ok := repoEntries[repositoryURL]
		if !ok {
			continue
		}

		// find entries for this chart
		if _, ok = entries[helmRelease.ChartName]; !ok {
			panic(fmt.Sprintf("No chart for name %s", helmRelease.ChartName))
		}

		foundVersions := []string{}
		for _, chartVersion := range entries[helmRelease.ChartName] {
			if _, ok = chartVersion.Annotations["artifacthub.io/prerelease"]; ok {
				if chartVersion.Annotations["artifacthub.io/prerelease"] == "true" {
					continue
				}
			}

			pattern := ".+"
			if _, ok = versions.Patterns.HelmCharts[chartName]; ok {
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

		if len(semvers) == 0 {
			continue
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

// GetImageUpdates checks Docker images for newer versions.
// Returns a map from image name to "currentVersion;latestVersion".
func (builder *Builder) GetImageUpdates(debug bool, filter string) (map[string]string, error) {
	versions := builder.Versions

	// Deduplicate images and build the work list.
	type imageWork struct {
		name    string
		version string
		pattern string
	}

	seen := make(map[string]bool)

	var work []imageWork

	for _, image := range builder.DockerImages {
		parts := strings.Split(image, ":")

		if seen[parts[0]] {
			continue
		}

		seen[parts[0]] = true

		if filter != "" && !strings.Contains(image, filter) {
			if debug {
				log.Printf("skipping image check: %s...", image)
			}

			continue
		}

		pattern := ".+"
		if p, ok := versions.Patterns.Images[parts[0]]; ok {
			pattern = p
		}

		work = append(work, imageWork{name: parts[0], version: parts[1], pattern: pattern})
	}

	// Check all images concurrently.
	var mu sync.Mutex

	retVersions := make(map[string]string)

	var wg sync.WaitGroup

	for _, w := range work {
		wg.Add(1)

		go func(item imageWork) {
			defer wg.Done()

			if debug {
				log.Printf("checking %s...", item.name)
			}

			startTs := time.Now()

			foundVersions := GetLastImageTag(debug, item.name, item.version, item.pattern)

			if len(foundVersions) > 0 {
				mu.Lock()
				retVersions[item.name] = fmt.Sprintf("%s;%s", item.version, foundVersions[len(foundVersions)-1])
				mu.Unlock()
			}

			if debug {
				log.Printf("done checking %s after %d msec", item.name, time.Since(startTs).Milliseconds())
			}
		}(w)
	}

	wg.Wait()

	return retVersions, nil
}

// CheckVersions checks the installed releases for update
func (builder *Builder) CheckVersions(debug bool, filter string) {
	var (
		helmVersions  map[string]string
		imageVersions map[string]string
		helmErr       error
		imageErr      error
	)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		helmVersions, helmErr = builder.GetHelmUpdates(debug, filter)
	}()

	go func() {
		defer wg.Done()
		imageVersions, imageErr = builder.GetImageUpdates(debug, filter)
	}()

	wg.Wait()

	if helmErr != nil {
		panic(helmErr)
	}

	if imageErr != nil {
		panic(imageErr)
	}

	for k, v := range helmVersions {
		fmt.Printf("%s;%s\n", k, v)
	}

	for k, v := range imageVersions {
		parts := strings.SplitN(v, ";", 2)
		fmt.Printf("%s;%s;%s\n", k, parts[0], parts[1])
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

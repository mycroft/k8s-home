package gitea

import (
	"context"
	"fmt"
	"log"
	"strings"
)

// CreateVersionBumpPRs creates one pull request per outdated dependency in updates.
//
// updates is a map from dependency name to "oldVersion;newVersion".
// If dryRun is true, the function prints what it would do without making API calls.
func CreateVersionBumpPRs(
	ctx context.Context,
	client *Client,
	updates map[string]string,
	baseBranch string,
	dryRun bool,
) error {
	if len(updates) == 0 {
		fmt.Println("All versions are up to date.")

		return nil
	}

	var errCount int

	for name, versionPair := range updates {
		parts := strings.SplitN(versionPair, ";", 2)
		oldVer, newVer := parts[0], parts[1]

		branchName := buildBranchName(name, oldVer)
		title := fmt.Sprintf("Updating %s from %s to %s", name, oldVer, newVer)

		fmt.Printf("Processing %s: %s → %s\n", name, oldVer, newVer)

		if dryRun {
			fmt.Printf("  [dry-run] would create branch %q and PR %q\n", branchName, title)

			continue
		}

		if err := createSinglePR(ctx, client, name, oldVer, newVer, branchName, baseBranch, title); err != nil {
			log.Printf("  error: %v", err)
			errCount++
		}
	}

	if errCount > 0 {
		return fmt.Errorf("%d PR(s) failed to create", errCount)
	}

	return nil
}

func buildBranchName(name, oldVer string) string {
	replacer := strings.NewReplacer("/", "_", ":", "_")

	return replacer.Replace(name + "_" + oldVer)
}

func createSinglePR(
	ctx context.Context,
	client *Client,
	name, oldVer, newVer, branchName, baseBranch, title string,
) error {
	exists, err := client.BranchExists(ctx, branchName)
	if err != nil {
		return fmt.Errorf("check branch: %w", err)
	}

	if exists {
		fmt.Printf("  branch %q already exists, skipping\n", branchName)

		return nil
	}

	if err = client.CreateBranch(ctx, branchName, baseBranch); err != nil {
		return fmt.Errorf("create branch: %w", err)
	}

	content, sha, err := client.GetFile(ctx, "versions.yaml", branchName)
	if err != nil {
		return fmt.Errorf("get versions.yaml: %w", err)
	}

	updated := strings.Replace(
		content,
		fmt.Sprintf("  %s: %s", name, oldVer),
		fmt.Sprintf("  %s: %s", name, newVer),
		1,
	)

	if updated == content {
		return fmt.Errorf("versions.yaml unchanged after substitution for %s", name)
	}

	commitMsg := fmt.Sprintf("Updating %s from %s to %s", name, oldVer, newVer)

	if err = client.UpdateFile(ctx, "versions.yaml", branchName, sha, commitMsg, updated); err != nil {
		return fmt.Errorf("update versions.yaml: %w", err)
	}

	prBody := fmt.Sprintf("Automated version bump: update `%s` from `%s` to `%s`.", name, oldVer, newVer)

	pr, err := client.CreatePR(ctx, title, branchName, baseBranch, prBody)
	if err != nil {
		return fmt.Errorf("create PR: %w", err)
	}

	fmt.Printf("  created PR #%d: %s\n", pr.Number, pr.Title)

	return nil
}

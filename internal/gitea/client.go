package gitea

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client is a minimal Gitea API client for repository operations.
type Client struct {
	baseURL string
	token   string
	owner   string
	repo    string
	http    *http.Client
}

// NewClient returns a new Gitea API client.
func NewClient(baseURL, token, owner, repo string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		owner:   owner,
		repo:    repo,
		http:    &http.Client{},
	}
}

// PRBranch holds the branch reference within a PullRequest.
type PRBranch struct {
	Label string `json:"label"`
}

// PullRequest is a simplified representation of a Gitea pull request.
type PullRequest struct {
	Number int      `json:"number"`
	Title  string   `json:"title"`
	State  string   `json:"state"`
	Head   PRBranch `json:"head"`
}

type fileContent struct {
	Content string `json:"content"`
	SHA     string `json:"sha"`
}

type updateFileRequest struct {
	Message string `json:"message"`
	Content string `json:"content"`
	SHA     string `json:"sha"`
	Branch  string `json:"branch"`
	Signoff bool   `json:"signoff"`
}

type createBranchRequest struct {
	NewBranchName string `json:"new_branch_name"`
	OldBranchName string `json:"old_branch_name"`
}

type createPRRequest struct {
	Title string `json:"title"`
	Head  string `json:"head"`
	Base  string `json:"base"`
	Body  string `json:"body"`
}

func apiError(action string, status int, body []byte) error {
	return fmt.Errorf("%s: unexpected status %d: %s", action, status, strings.TrimSpace(string(body)))
}

func (c *Client) doRequest(ctx context.Context, method, path string, payload any) ([]byte, int, error) {
	var bodyReader io.Reader

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal payload: %w", err)
		}

		bodyReader = bytes.NewReader(b)
	}

	url := c.baseURL + "/api/v1" + path

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read body: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

// BranchExists returns true if the named branch exists in the repository.
func (c *Client) BranchExists(ctx context.Context, branch string) (bool, error) {
	path := fmt.Sprintf("/repos/%s/%s/branches/%s", c.owner, c.repo, branch)

	_, status, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return false, err
	}

	return status == http.StatusOK, nil
}

// CreateBranch creates a new branch from baseBranch.
func (c *Client) CreateBranch(ctx context.Context, newBranch, baseBranch string) error {
	path := fmt.Sprintf("/repos/%s/%s/branches", c.owner, c.repo)
	payload := createBranchRequest{
		NewBranchName: newBranch,
		OldBranchName: baseBranch,
	}

	body, status, err := c.doRequest(ctx, http.MethodPost, path, payload)
	if err != nil {
		return err
	}

	if status != http.StatusCreated {
		return apiError(fmt.Sprintf("create branch %q", newBranch), status, body)
	}

	return nil
}

// GetFile retrieves the decoded content and blob SHA of a file on a given branch.
func (c *Client) GetFile(ctx context.Context, filepath, branch string) (string, string, error) {
	path := fmt.Sprintf("/repos/%s/%s/contents/%s?ref=%s", c.owner, c.repo, filepath, branch)

	body, status, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return "", "", err
	}

	if status != http.StatusOK {
		return "", "", apiError(fmt.Sprintf("get file %q", filepath), status, body)
	}

	var fc fileContent
	if err = json.Unmarshal(body, &fc); err != nil {
		return "", "", fmt.Errorf("unmarshal file content: %w", err)
	}

	// Gitea encodes content as base64 with embedded newlines; strip them before decoding.
	rawContent := strings.ReplaceAll(fc.Content, "\n", "")

	decoded, err := base64.StdEncoding.DecodeString(rawContent)
	if err != nil {
		return "", "", fmt.Errorf("decode file content: %w", err)
	}

	return string(decoded), fc.SHA, nil
}

// UpdateFile commits an updated version of a file to a branch.
func (c *Client) UpdateFile(ctx context.Context, filepath, branch, sha, message, content string) error {
	path := fmt.Sprintf("/repos/%s/%s/contents/%s", c.owner, c.repo, filepath)
	payload := updateFileRequest{
		Message: message,
		Content: base64.StdEncoding.EncodeToString([]byte(content)),
		SHA:     sha,
		Branch:  branch,
		Signoff: true,
	}

	body, status, err := c.doRequest(ctx, http.MethodPut, path, payload)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return apiError(fmt.Sprintf("update file %q", filepath), status, body)
	}

	return nil
}

// ListOpenPRs returns all open pull requests in the repository.
func (c *Client) ListOpenPRs(ctx context.Context) ([]PullRequest, error) {
	path := fmt.Sprintf("/repos/%s/%s/pulls?state=open&limit=50", c.owner, c.repo)

	body, status, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, apiError("list PRs", status, body)
	}

	var prs []PullRequest
	if err = json.Unmarshal(body, &prs); err != nil {
		return nil, fmt.Errorf("unmarshal PR list: %w", err)
	}

	return prs, nil
}

type mergePRRequest struct {
	Do string `json:"Do"`
}

// MergePR merges a pull request by number using rebase strategy.
func (c *Client) MergePR(ctx context.Context, number int) error {
	path := fmt.Sprintf("/repos/%s/%s/pulls/%d/merge", c.owner, c.repo, number)
	payload := mergePRRequest{Do: "rebase"}

	body, status, err := c.doRequest(ctx, http.MethodPost, path, payload)
	if err != nil {
		return err
	}

	if status != http.StatusNoContent && status != http.StatusOK {
		return apiError(fmt.Sprintf("merge PR #%d", number), status, body)
	}

	return nil
}

// CreatePR opens a new pull request and returns it.
func (c *Client) CreatePR(ctx context.Context, title, head, base, body string) (*PullRequest, error) {
	path := fmt.Sprintf("/repos/%s/%s/pulls", c.owner, c.repo)
	payload := createPRRequest{
		Title: title,
		Head:  head,
		Base:  base,
		Body:  body,
	}

	respBody, status, err := c.doRequest(ctx, http.MethodPost, path, payload)
	if err != nil {
		return nil, err
	}

	if status != http.StatusCreated {
		return nil, apiError(fmt.Sprintf("create PR %q", title), status, respBody)
	}

	var pr PullRequest
	if err = json.Unmarshal(respBody, &pr); err != nil {
		return nil, fmt.Errorf("unmarshal PR response: %w", err)
	}

	return &pr, nil
}

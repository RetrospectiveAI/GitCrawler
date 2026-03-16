package service

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

type CloneService struct {
}

func NewCloneService() *CloneService {
	return &CloneService{}
}

// CloneRepository clones repositoryUrl into a fresh temp directory.
// When token is non-empty it is injected as the password in the HTTPS URL so
// that private GitHub / GitLab repositories can be cloned without interactive
// credentials.  The token is never written to logs.
func (c *CloneService) CloneRepository(repositoryUrl, token string) (string, error) {
	path, err := c.createRepositoryDirectory()
	if err != nil {
		return path, err
	}

	cloneUrl := repositoryUrl
	if token != "" {
		cloneUrl = c.injectToken(repositoryUrl, token)
	}

	cmd := exec.Command("git", "clone", cloneUrl, ".")
	cmd.Dir = path
	err = cmd.Run()
	if err != nil {
		return path, errors.New(fmt.Sprintf("Repository not found, project may be private: %s", err.Error()))
	}
	return path, nil
}

// injectToken rewrites an HTTPS GitHub/GitLab URL to embed the PAT so git
// can authenticate without a credential helper.
// Input:  https://github.com/owner/repo.git
// Output: https://x-access-token:<TOKEN>@github.com/owner/repo.git
// The raw token is never surfaced in error messages or log output.
func (c *CloneService) injectToken(rawUrl, token string) string {
	// Only rewrite HTTPS URLs
	if !strings.HasPrefix(rawUrl, "https://") {
		return rawUrl
	}
	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return rawUrl
	}
	parsed.User = url.UserPassword("x-access-token", token)
	return parsed.String()
}

func (c *CloneService) createRepositoryDirectory() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path, err = os.MkdirTemp(path, "temp")
	return path, nil
}

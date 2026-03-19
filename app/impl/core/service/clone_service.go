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

func (c *CloneService) CloneRepository(repositoryUrl, token string) (string, error) {
	c.normalizeGitUrl(repositoryUrl)
	path, err := c.createRepositoryDirectory()
	if err != nil {
		return path, err
	}

	cloneUrl := repositoryUrl
	if token != "" {
		cloneUrl = c.injectGitToken(repositoryUrl, token)
	}

	cmd := exec.Command("git", "clone", cloneUrl, ".")
	cmd.Dir = path
	err = cmd.Run()
	if err != nil {
		return path, errors.New(fmt.Sprintf("Repository not found, project may be private: %s", err.Error()))
	}
	return path, nil
}

func (c *CloneService) injectGitToken(rawUrl, token string) string {
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
func (c *CloneService) normalizeGitUrl(url string) string {
	url = strings.TrimRight(url, "/")
	if (strings.Contains(url, "github.com") || strings.Contains(url, "gitlab.com")) &&
		!strings.HasSuffix(url, ".git") {
		return url + ".git"
	}
	return url
}

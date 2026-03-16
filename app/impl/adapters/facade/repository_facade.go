package facade

import (
	"errors"
	contract2 "gitcrawler/app/impl/core/contract"
	"gitcrawler/app/impl/core/entity"
	"gitcrawler/app/impl/core/enum"
	strategy2 "gitcrawler/app/impl/core/strategy"
	"os"
	"path/filepath"
	"strings"
)

type RepositoryFacade struct {
	cloneService           contract2.CloneServiceContract
	crawlerService         contract2.CrawlerServiceContract
	resumeGeneratorService contract2.ResumeGeneratorServiceContract
}

func NewRepositoryFacade(cloneService contract2.CloneServiceContract, crawlerService contract2.CrawlerServiceContract, resumeGeneratorService contract2.ResumeGeneratorServiceContract) *RepositoryFacade {
	return &RepositoryFacade{
		cloneService:           cloneService,
		crawlerService:         crawlerService,
		resumeGeneratorService: resumeGeneratorService,
	}
}

func (c *RepositoryFacade) GetRepositoryFiles(url string, extensions []string, dirs []string, token string) (data *entity.RepositoryData, err error) {
	url = c.normalizeUrl(url)
	err = c.isUrlValid(url)
	if err != nil {
		return nil, err
	}
	data, err = c.createAndCrawl(url, extensions, dirs, token)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func (c *RepositoryFacade) SaveRepositoryFiles(url string, extensions []string, dirs []string, option enum.ConversionOption) (err error) {
	url = c.normalizeUrl(url)
	err = c.isUrlValid(url)
	if err != nil {
		return err
	}
	data, err := c.createAndCrawl(url, extensions, dirs, "")

	converter, err := c.converterStrategy(option)

	if err != nil {
		return err
	}
	err = converter.Convert(data)

	if err != nil {
		return err
	}
	return nil
}

func (c *RepositoryFacade) GenerateBusinessResume(url, token string) (aiResponse string, err error) {
	extensions := []string{
		".java",
		".kt",
		".kts",
		".groovy",
		".go",
		".py",
		".ts",
		".js",
		".cs",
		".rb",
		".php",
		".rs",
		".cpp",
		".cc",
		".cxx",
		".c",
		".h",
		".hpp",
		".scala",
		".ex",
		".exs",
		".dart",
		".swift",
		".md",
	}

	// Focused dirs for business-logic detection (first pass)
	focusedDirs := []string{
		"rest",
		"api",
		"controller",
		"handler",
		"service",
		"usecase",
		"application",
		"domain",
		"model",
		"entity",
		"aggregate",
		"facade",
		"client",
		"gateway",
		"integration",
		"repository",
		"event",
		"main",
	}

	// Broader fallback dirs that cover flat / non-standard layouts
	broadDirs := []string{
		"src", "lib", "pkg", "app", "core", "internal", "cmd",
		"main", "server", "backend", "frontend", "common",
		"utils", "helpers", "middleware", "routes", "handlers",
		"modules", "components", "features", "pages",
		"rest", "api", "controller", "handler",
		"service", "usecase", "application", "domain",
		"model", "entity", "aggregate", "facade",
		"client", "gateway", "integration", "repository", "event",
	}

	url = c.normalizeUrl(url)
	err = c.isUrlValid(url)
	if err != nil {
		return "", err
	}

	data, err := c.createAndCrawl(url, extensions, focusedDirs, token)
	if err != nil {
		return "", err
	}
	if data == nil || len(data.Files) == 0 {
		// First pass found nothing – retry with the broader directory list
		data, err = c.createAndCrawl(url, extensions, broadDirs, token)
		if err != nil {
			return "", err
		}
	}
	if data == nil || len(data.Files) == 0 {
		return "", errors.New("no source files found in repository (check that the URL is correct and the repo is public)")
	}
	aiResponse, err = c.resumeGeneratorService.GenerateBusinessResume(data.String())
	if err != nil {
		return "", err
	}
	return aiResponse, nil
}

func (c *RepositoryFacade) createAndCrawl(url string, extensions []string, dirs []string, token string) (data *entity.RepositoryData, err error) {
	path, err := c.cloneService.CloneRepository(url, token)
	if path != "" {
		defer os.RemoveAll(path)
	}

	if err != nil {
		return nil, err
	}
	repoName := strings.TrimSuffix(filepath.Base(url), ".git")

	data, err = c.crawlerService.CrawlRepository(path, repoName, extensions, dirs)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *RepositoryFacade) converterStrategy(option enum.ConversionOption) (converter strategy2.DataConverter, err error) {
	switch option {
	case enum.Csv:
		return strategy2.NewConverterCsv(), nil
	default:
		return nil, errors.New("unknown conversion option")
	}
}

// normalizeUrl appends ".git" to GitHub/GitLab URLs that are missing it,
// so users can paste plain browser URLs without getting a validation error.
func (c *RepositoryFacade) normalizeUrl(url string) string {
	url = strings.TrimRight(url, "/")
	if (strings.Contains(url, "github.com") || strings.Contains(url, "gitlab.com")) &&
		!strings.HasSuffix(url, ".git") {
		return url + ".git"
	}
	return url
}

func (c *RepositoryFacade) isUrlValid(url string) error {
	if !strings.Contains(url, ".git") {
		return errors.New("repository url must contain .git – pass a valid GitHub/GitLab clone URL")
	}
	return nil
}

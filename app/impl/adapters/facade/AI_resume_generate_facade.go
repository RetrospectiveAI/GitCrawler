package facade

import (
	"errors"
	"gitcrawler/app/impl/core/contract"
	"gitcrawler/app/impl/core/service"
)

type AIResumeGenerateFacade struct {
	repositoryLoader       *service.RepositoryLoaderService
	resumeGeneratorService contract.ResumeGeneratorServiceContract
}

func NewAIResumeGenerateFacade(loaderService *service.RepositoryLoaderService, resumeGeneratorService contract.ResumeGeneratorServiceContract) *AIResumeGenerateFacade {
	return &AIResumeGenerateFacade{
		repositoryLoader:       loaderService,
		resumeGeneratorService: resumeGeneratorService,
	}

}

func (c *AIResumeGenerateFacade) GenerateBusinessResume(url, token string) (aiResponse string, err error) {
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

	fallbackDirs := []string{
		"src", "lib", "pkg", "app", "core", "internal", "cmd",
		"main", "server", "backend", "frontend", "common",
		"utils", "helpers", "middleware", "routes", "handlers",
		"modules", "components", "features", "pages",
		"rest", "api", "controller", "handler",
		"service", "usecase", "application", "domain",
		"model", "entity", "aggregate", "facade",
		"client", "gateway", "integration", "repository", "event",
	}

	data, err := c.repositoryLoader.CreateAndCrawl(url, extensions, focusedDirs, token)
	if err != nil {
		return "", err
	}
	if data == nil || len(data.Files) == 0 {
		data, err = c.repositoryLoader.CreateAndCrawl(url, extensions, fallbackDirs, token)
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

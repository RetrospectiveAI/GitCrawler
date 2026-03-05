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

func (c *RepositoryFacade) GetRepositoryFiles(url string, extensions []string, dirs []string) (data *entity.RepositoryData, err error) {
	err = c.isUrlValid(url)
	if err != nil {
		return nil, err
	}
	data, err = c.createAndCrawl(url, extensions, dirs)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func (c *RepositoryFacade) SaveRepositoryFiles(url string, extensions []string, dirs []string, option enum.ConversionOption) (err error) {
	err = c.isUrlValid(url)
	if err != nil {
		return err
	}
	data, err := c.createAndCrawl(url, extensions, dirs)

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

func (c *RepositoryFacade) GenerateBusinessResume(url string) (aiResponse string, err error) {
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

	dirs := []string{
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
	}
	err = c.isUrlValid(url)
	if err != nil {
		return "", err
	}

	data, err := c.createAndCrawl(url, extensions, dirs)
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", errors.New("repository business data is empty")
	}
	if len(data.Files) < 5 {
		return "", errors.New("repository business context is too small")
	}
	aiResponse, err = c.resumeGeneratorService.GenerateBusinessResume(data.String())
	if err != nil {
		return "", err
	}
	return aiResponse, nil
}

func (c *RepositoryFacade) createAndCrawl(url string, extensions []string, dirs []string) (data *entity.RepositoryData, err error) {
	path, err := c.cloneService.CloneRepository(url)
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

func (c *RepositoryFacade) isUrlValid(url string) error {
	if !strings.Contains(url, ".git") {
		return errors.New("repository url must contain .git")
	}
	return nil
}

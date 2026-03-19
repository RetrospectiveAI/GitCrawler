package service

import (
	"gitcrawler/app/impl/core/contract"
	"gitcrawler/app/impl/core/model"
	"os"
	"path/filepath"
	"strings"
)

type RepositoryLoaderService struct {
	crawlerService contract.CrawlerServiceContract
	cloneService   contract.CloneServiceContract
}

func NewRepositoryLoaderService(crawlerService contract.CrawlerServiceContract, cloneService contract.CloneServiceContract) *RepositoryLoaderService {
	return &RepositoryLoaderService{
		cloneService:   cloneService,
		crawlerService: crawlerService,
	}

}

func (c *RepositoryLoaderService) CreateAndCrawl(url string, extensions []string, dirs []string, token string) (data *model.RepositoryData, err error) {
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

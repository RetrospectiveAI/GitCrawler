package contract

import (
	"gitcrawler/app/impl/core/model"
)

type CrawlerServiceContract interface {
	CrawlRepository(path string, repoName string, validExtensions []string, validDirs []string) (data *model.RepositoryData, err error)
}

package facade

import (
	"errors"
	"gitcrawler/app/impl/core/contract"
	"gitcrawler/app/impl/core/enum"
	"gitcrawler/app/impl/core/model"
	"gitcrawler/app/impl/core/service"
	"gitcrawler/app/impl/core/strategy"
	"path/filepath"
	"strings"
)

type RepositoryFacade struct {
	repositoryLoader *service.RepositoryLoaderService
	fileWriter       contract.FileWriterContract
}

func NewRepositoryFacade(repositoryLoader *service.RepositoryLoaderService, fileWriter contract.FileWriterContract) *RepositoryFacade {
	return &RepositoryFacade{
		repositoryLoader: repositoryLoader,
		fileWriter:       fileWriter,
	}
}

func (c *RepositoryFacade) GetRepositoryFiles(url string, extensions []string, dirs []string, token string) (data *model.RepositoryData, err error) {
	data, err = c.repositoryLoader.CreateAndCrawl(url, extensions, dirs, token)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func (c *RepositoryFacade) SaveRepositoryFiles(url string, extensions []string, dirs []string, option enum.ConversionOption, token string) (err error) {
	data, err := c.repositoryLoader.CreateAndCrawl(url, extensions, dirs, token)

	if err != nil {
		return err
	}

	converter, err := c.converterStrategy(option)

	if err != nil {
		return err
	}
	repositoryData, err := converter.Convert(data)
	if err != nil {
		return err
	}
	repoName := strings.TrimSuffix(filepath.Base(url), ".git")

	err = c.fileWriter.WriteConvertedFiles(repositoryData, repoName, string(option))

	if err != nil {
		return err
	}

	return nil
}

func (c *RepositoryFacade) converterStrategy(option enum.ConversionOption) (converter strategy.DataConverter, err error) {
	switch option {
	case enum.Csv:
		return strategy.NewConverterCsv(), nil

	case enum.Json:
		return strategy.NewConverterJson(), nil
	default:
		return nil, errors.New("unknown conversion option")
	}
}

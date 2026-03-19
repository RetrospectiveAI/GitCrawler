package service

import (
	"errors"
	"gitcrawler/app/impl/core/model"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type CrawlerService struct {
	extensionList []string
	dirList       []string
}

func NewCrawlerService() *CrawlerService {
	return &CrawlerService{}
}

func (c *CrawlerService) CrawlRepository(path string, repoName string, validExtensions []string, validDirs []string) (data *model.RepositoryData, err error) {
	data = &model.RepositoryData{}
	data.Name = repoName
	if validExtensions == nil {
		return nil, errors.New("no valid extensions provided")
	}
	if validDirs == nil {
		return nil, errors.New("no valid directories provided")
	}
	c.extensionList = validExtensions
	c.dirList = validDirs

	err = c.crawl(path, data, true)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (c *CrawlerService) crawl(dir string, repositoryData *model.RepositoryData, validDir bool) (err error) {
	readDir, err := os.ReadDir(dir)

	if err != nil {
		return err
	}
	for _, file := range readDir {
		path := filepath.Join(dir, file.Name())
		isDirValid := c.isDirValid(file)
		if file.IsDir() {
			err = c.crawl(path, repositoryData, validDir || isDirValid)
			if err != nil {
				return err
			}
		} else {
			if c.isFileValid(file) && validDir {
				fileData, openFileErr := c.openFile(path)
				if openFileErr != nil {
					return openFileErr
				}
				c.appendFileData(repositoryData, path, fileData)
			}
		}
	}
	return nil
}

func (c *CrawlerService) openFile(path string) (data string, err error) {
	fileBinary, err := os.Open(path)

	if err != nil {
		return "", errors.New("could not open file")
	}
	defer fileBinary.Close()
	binaryData, err := io.ReadAll(fileBinary)

	if err != nil {
		return "", errors.New("could not read file")
	}
	data = string(binaryData)
	return data, nil
}

func (c *CrawlerService) appendFileData(repositoryData *model.RepositoryData, path string, data string) {
	repositoryFile := &model.RepositoryFile{}
	repositoryFile.Data = data
	repositoryFile.Path = path
	repositoryData.Files = append(repositoryData.Files, repositoryFile)
}
func (c *CrawlerService) isFileValid(file os.DirEntry) bool {
	fileExtension := strings.ToLower(filepath.Ext(file.Name()))
	for _, extension := range c.extensionList {
		if fileExtension == extension {
			return true
		}
	}
	return false
}
func (c *CrawlerService) isDirValid(file os.DirEntry) bool {
	dirName := strings.ToLower(file.Name())
	for _, dir := range c.dirList {
		if dir == dirName {
			return true
		}
	}
	return false
}

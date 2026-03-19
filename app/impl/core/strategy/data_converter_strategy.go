package strategy

import (
	"gitcrawler/app/impl/core/model"
)

type DataConverter interface {
	Convert(data *model.RepositoryData) (value []byte, err error)
}

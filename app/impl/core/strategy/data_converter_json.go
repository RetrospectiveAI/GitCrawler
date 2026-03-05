package strategy

import (
	"encoding/json"
	"gitcrawler/app/impl/core/entity"
)

type ConverterJson struct{}

func NewConverterJson() *ConverterJson {
	return &ConverterJson{}
}

func (c *ConverterJson) Convert(data *entity.RepositoryData) ([]byte, error) {
	return json.Marshal(data)
}

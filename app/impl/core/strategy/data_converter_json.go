package strategy

import (
	"encoding/json"
	"gitcrawler/app/impl/core/model"
)

type ConverterJson struct{}

func NewConverterJson() *ConverterJson {
	return &ConverterJson{}
}

func (c *ConverterJson) Convert(data *model.RepositoryData) ([]byte, error) {
	return json.Marshal(data)
}

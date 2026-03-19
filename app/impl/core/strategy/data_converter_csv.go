package strategy

import (
	"bytes"
	"encoding/csv"
	"gitcrawler/app/impl/core/model"
	"strings"
)

type ConverterCsv struct{}

func NewConverterCsv() *ConverterCsv {
	return &ConverterCsv{}
}

func (c *ConverterCsv) Convert(data *model.RepositoryData) (value []byte, err error) {

	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)

	headers := []string{"Name", "Data", "Path"}

	writer.Write(headers)

	for i := 0; i < len(data.Files); i++ {
		trimmedData := strings.ReplaceAll(data.Files[i].Data, "\n", "\\n")
		record := []string{
			data.Name,
			trimmedData,
			data.Files[i].Path,
		}
		err = writer.Write(record)
		if err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return buffer.Bytes(), nil
}

package integration

import (
	"os"
	"path/filepath"
)

type FileWriter struct {
}

func NewFileWriter() *FileWriter {
	return &FileWriter{}
}
func (s *FileWriter) WriteConvertedFiles(data []byte, repositoryName string, extension string) (err error) {
	homeDir, _ := os.UserHomeDir()

	dirPath := homeDir + "/Downloads"

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := repositoryName + "." + extension
	filePath := filepath.Join(dirPath, fileName)

	writer, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	return nil
}

package contract

type FileWriterContract interface {
	WriteConvertedFiles(data []byte, repositoryName string, extension string) (err error)
}

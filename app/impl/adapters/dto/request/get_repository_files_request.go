package request

type GetRepositoryFilesRequest struct {
	Url        string   `json:"url"`
	Dirs       []string `json:"dirs"`
	Extensions []string `json:"extensions"`
}

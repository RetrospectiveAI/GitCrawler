package response

type RepositoryDataResponse struct {
	Files []*RepositoryFileResponse `json:"files"`
	Name  string                    `json:"name"`
}

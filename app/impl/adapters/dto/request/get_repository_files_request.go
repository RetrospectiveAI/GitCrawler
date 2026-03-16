package request

type GetRepositoryFilesRequest struct {
	Url        string   `json:"url"`
	Dirs       []string `json:"dirs"`
	Extensions []string `json:"extensions"`
	// Token is an optional GitHub/GitLab PAT for private repositories.
	// It is intentionally excluded from any JSON responses or log output.
	Token string `json:"token,omitempty"`
}

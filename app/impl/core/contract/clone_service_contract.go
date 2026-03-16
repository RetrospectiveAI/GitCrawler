package contract

type CloneServiceContract interface {
	// CloneRepository clones repositoryUrl into a temp directory.
	// token is a GitHub/GitLab Personal Access Token; pass "" for public repos.
	CloneRepository(repositoryUrl, token string) (string, error)
}

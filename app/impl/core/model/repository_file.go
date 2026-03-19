package model

import "fmt"

type RepositoryFile struct {
	Data string
	Path string
}

func (r *RepositoryFile) String() string {
	return fmt.Sprintf("=== FILE: %s ===\n%s", r.Path, r.Data)
}

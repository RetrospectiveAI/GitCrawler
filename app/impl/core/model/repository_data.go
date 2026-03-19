package model

import (
	"fmt"
	"strings"
)

type RepositoryData struct {
	Files []*RepositoryFile
	Name  string
}

// MaxResumeBytes caps the content fed to the LLM for resume generation.
// ~80 KB sits comfortably within 128k-token models while still covering the
// most important parts of any realistic codebase.  Files beyond this limit
// are listed in the file tree but their content is omitted.
const MaxResumeBytes = 80_000

// String produces a structured, LLM-readable representation of the repository.
//
// Format:
//
//	=== REPOSITORY: <name> ===
//	FILE TREE:
//	  src/main/Foo.java
//	  ...
//	=== FILE: src/main/Foo.java ===
//	<content>
//
// The full file tree is always included so the LLM knows every path that
// exists.  File contents are included in order until MaxResumeBytes is
// reached; remaining files are noted but their bodies are omitted.
func (r *RepositoryData) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "=== REPOSITORY: %s ===\n\n", r.Name)

	// Always emit the complete file tree regardless of size limit
	sb.WriteString("FILE TREE:\n")
	for _, f := range r.Files {
		fmt.Fprintf(&sb, "  %s\n", f.Path)
	}
	sb.WriteString("\n")

	// Emit file contents up to the size cap
	omitted := 0
	for _, f := range r.Files {
		if sb.Len() >= MaxResumeBytes {
			omitted++
			continue
		}
		fmt.Fprintf(&sb, "=== FILE: %s ===\n%s\n\n", f.Path, f.Data)
	}
	if omitted > 0 {
		fmt.Fprintf(&sb, "[%d additional file(s) omitted – content exceeds context budget]\n", omitted)
	}

	return sb.String()
}

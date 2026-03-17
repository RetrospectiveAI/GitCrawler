package request

import (
	"gitcrawler/app/impl/core/enum"
)

type SaveRepositoryFilesRequest struct {
	Url        string                `json:"url"`
	Dirs       []string              `json:"dirs"`
	Extensions []string              `json:"extensions"`
	Option     enum.ConversionOption `json:"option"`
	Token      string                `json:"token,omitempty"`
}

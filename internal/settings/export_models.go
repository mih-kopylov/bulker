package settings

import "golang.org/x/exp/slices"

type ExportImportStatus int

const (
	// ExportImportStatusCompleted means the operation completed successfully
	ExportImportStatusCompleted ExportImportStatus = iota
	// ExportImportStatusUpToDate means there was no change
	ExportImportStatusUpToDate
)

type exportModel struct {
	Version int         `yaml:"version"`
	Data    modelDataV1 `yaml:"data"`
}

type modelDataV1 struct {
	Repos map[string]modelDataV1Repo `yaml:"repos"`
}

type modelDataV1Repo struct {
	Url  string   `yaml:"url"`
	Tags []string `yaml:"tags"`
}

func (r modelDataV1Repo) Equals(other modelDataV1Repo) bool {
	return r.Url == other.Url && slices.Equal(r.Tags, other.Tags)
}

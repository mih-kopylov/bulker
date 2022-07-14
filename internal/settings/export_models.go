package settings

import "golang.org/x/exp/slices"

type ExportStatus int

const (
	ExportStatusExported ExportStatus = iota
	ExportStatusUpToDate
)

type exportModel struct {
	Version int               `yaml:"version"`
	Data    exportModelDataV1 `yaml:"data"`
}

type importModel struct {
	Version int `yaml:"version"`
}

type exportModelDataV1 struct {
	Repos map[string]exportModelDataV1Repo `yaml:"repos"`
}

type exportModelDataV1Repo struct {
	Url  string   `yaml:"url"`
	Tags []string `yaml:"tags"`
}

func (r exportModelDataV1Repo) Equals(other exportModelDataV1Repo) bool {
	return r.Url == other.Url && slices.Equal(r.Tags, other.Tags)
}

func fromSettings(settings *Settings) *exportModel {
	data := exportModelDataV1{}
	data.Repos = map[string]exportModelDataV1Repo{}
	for _, repo := range settings.Repos {
		data.Repos[repo.Name] = exportModelDataV1Repo{
			Url:  repo.Url,
			Tags: repo.Tags,
		}
	}

	return &exportModel{1, data}
}

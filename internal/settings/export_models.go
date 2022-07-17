package settings

import "golang.org/x/exp/slices"

type ExportStatus int

const (
	ExportStatusExported ExportStatus = iota
	ExportStatusUpToDate
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

func fromSettings(settings *Settings) *exportModel {
	data := modelDataV1{}
	data.Repos = map[string]modelDataV1Repo{}
	for _, repo := range settings.Repos {
		data.Repos[repo.Name] = modelDataV1Repo{
			Url:  repo.Url,
			Tags: repo.Tags,
		}
	}

	return &exportModel{1, data}
}

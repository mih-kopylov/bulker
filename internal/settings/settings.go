package settings

import (
	"fmt"
	"sort"
	"strings"
)

type Settings struct {
	Repos []Repo `yaml:"repos"`
}

type Repo struct {
	Name string   `yaml:"name"`
	Url  string   `yaml:"url"`
	Tags []string `yaml:"tags"`
}

func (s *Settings) AddRepo(name string, url string, tags []string) error {
	repoIndex := s.findRepoIndex(name)

	if repoIndex >= 0 {
		return fmt.Errorf("repository %v already exists", name)
	}

	repo := Repo{
		Name: name,
		Url:  url,
		Tags: tags,
	}
	s.Repos = append(s.Repos, repo)

	// make sure repos are sorted alphabetically
	sort.Slice(
		s.Repos, func(i, j int) bool {
			return strings.Compare(s.Repos[i].Name, s.Repos[j].Name) > 0
		},
	)

	return nil
}

func (s *Settings) RemoveRepo(name string) error {
	repoIndex := s.findRepoIndex(name)

	if repoIndex < 0 {
		return fmt.Errorf("repository %v is not found", name)
	}

	s.Repos = append(s.Repos[:repoIndex], s.Repos[repoIndex+1:]...)

	return nil
}

func (s *Settings) findRepoIndex(name string) int {
	repoIndex := -1

	for i, repo := range s.Repos {
		if repo.Name == name {
			repoIndex = i
		}
	}

	return repoIndex
}

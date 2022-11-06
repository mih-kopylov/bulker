package settings

import (
	"errors"
	"golang.org/x/exp/slices"
)

type Settings struct {
	Repos  []Repo  `yaml:"repos"`
	Groups []Group `yaml:"groups"`
}

type Repo struct {
	Name string   `yaml:"name"`
	Url  string   `yaml:"url"`
	Tags []string `yaml:"tags"`
}

type Group struct {
	Name  string   `yaml:"name"`
	Repos []string `yaml:"repos"`
}

const (
	PreviousGroupName = "previous"
)

var (
	ErrGroupNotFound      = errors.New("group is not found")
	ErrGroupAlreadyExists = errors.New("group already exists")
	ErrRepoAlreadyExists  = errors.New("repository already exists")
	ErrRepoNotFound       = errors.New("repository is not found")
	ErrRepoNotSupported   = errors.New("repository is not supported")
	ErrRepoAlreadyAdded   = errors.New("repository already added")
	ErrRepoAlreadyRemoved = errors.New("repository already removed")
)

func (s *Settings) AddRepo(name string, url string, tags []string) error {
	if s.RepoExists(name) {
		return ErrRepoAlreadyExists
	}

	repo := Repo{
		Name: name,
		Url:  url,
		Tags: tags,
	}
	s.Repos = append(s.Repos, repo)

	return nil
}

func (s *Settings) RemoveRepo(name string) error {
	repoIndex := s.getRepoIndex(name)

	if repoIndex < 0 {
		return ErrRepoNotFound
	}

	s.Repos = slices.Delete(s.Repos, repoIndex, repoIndex+1)

	return nil
}

func (s *Settings) RepoExists(name string) bool {
	return s.getRepoIndex(name) >= 0
}

func (s *Settings) getRepoIndex(name string) int {
	return slices.IndexFunc(
		s.Repos, func(repo Repo) bool {
			return repo.Name == name
		},
	)
}

func (s *Settings) getGroupIndex(group string) int {
	return slices.IndexFunc(
		s.Groups, func(g Group) bool {
			return g.Name == group
		},
	)
}

func (s *Settings) GetGroup(group string) (*Group, error) {
	groupIndex := s.getGroupIndex(group)

	if groupIndex < 0 {
		return nil, ErrGroupNotFound
	}

	return &s.Groups[groupIndex], nil

}

func (s *Settings) GroupExists(group string) bool {
	return s.getGroupIndex(group) >= 0
}

func (s *Settings) RemoveGroup(group string) error {
	groupIndex := s.getGroupIndex(group)

	if groupIndex < 0 {
		return ErrGroupNotFound
	}

	s.Groups = slices.Delete(s.Groups, groupIndex, groupIndex+1)
	return nil
}

func (s *Settings) AddGroup(group string) error {
	if s.GroupExists(group) {
		return ErrGroupAlreadyExists
	}

	newGroup := Group{
		Name:  group,
		Repos: []string{},
	}

	s.Groups = append(s.Groups, newGroup)

	return nil
}

func (s *Settings) AddRepoToGroup(groupName string, repoName string) error {
	if !s.RepoExists(repoName) {
		return ErrRepoNotSupported
	}

	group, err := s.GetGroup(groupName)
	if err != nil {
		return err
	}

	if slices.Contains(group.Repos, repoName) {
		return ErrRepoAlreadyAdded
	}

	group.Repos = append(group.Repos, repoName)

	return nil
}

func (s *Settings) RemoveRepoFromGroup(groupName string, repoName string) error {
	if !s.RepoExists(repoName) {
		return ErrRepoNotSupported
	}

	group, err := s.GetGroup(groupName)
	if err != nil {
		return err
	}

	repoIndex := slices.Index(group.Repos, repoName)
	if repoIndex < 0 {
		return ErrRepoAlreadyRemoved
	}

	group.Repos = slices.Delete(group.Repos, repoIndex, repoIndex+1)
	return nil
}

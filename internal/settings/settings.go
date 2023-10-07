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

func (s *Settings) GetRepo(name string) (*Repo, error) {
	repoIndex := s.getRepoIndex(name)
	if repoIndex < 0 {
		return nil, ErrRepoNotFound
	}

	repo := s.Repos[repoIndex]
	return &repo, nil
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

func (s *Settings) getGroupIndex(groupName string) int {
	return slices.IndexFunc(
		s.Groups, func(g Group) bool {
			return g.Name == groupName
		},
	)
}

func (s *Settings) GetGroup(groupName string) (*Group, error) {
	groupIndex := s.getGroupIndex(groupName)

	if groupIndex < 0 {
		return nil, ErrGroupNotFound
	}

	return &s.Groups[groupIndex], nil

}

func (s *Settings) GroupExists(groupName string) bool {
	return s.getGroupIndex(groupName) >= 0
}

func (s *Settings) RemoveGroup(groupName string) error {
	groupIndex := s.getGroupIndex(groupName)

	if groupIndex < 0 {
		return ErrGroupNotFound
	}

	s.Groups = slices.Delete(s.Groups, groupIndex, groupIndex+1)
	return nil
}

func (s *Settings) AddGroup(groupName string) (*Group, error) {
	if s.GroupExists(groupName) {
		return nil, ErrGroupAlreadyExists
	}

	group := Group{
		Name:  groupName,
		Repos: []string{},
	}

	s.Groups = append(s.Groups, group)

	storedGroup, err := s.GetGroup(group.Name)
	if err != nil {
		return nil, err
	}

	return storedGroup, nil
}

func (s *Settings) AddRepoToGroup(group *Group, repoName string) error {
	if !s.RepoExists(repoName) {
		return ErrRepoNotSupported
	}

	if slices.Contains(group.Repos, repoName) {
		return ErrRepoAlreadyAdded
	}

	group.Repos = append(group.Repos, repoName)

	return nil
}

func (s *Settings) RemoveRepoFromGroup(group *Group, repoName string) error {
	if !s.RepoExists(repoName) {
		return ErrRepoNotSupported
	}

	repoIndex := slices.Index(group.Repos, repoName)
	if repoIndex < 0 {
		return ErrRepoAlreadyRemoved
	}

	group.Repos = slices.Delete(group.Repos, repoIndex, repoIndex+1)
	return nil
}

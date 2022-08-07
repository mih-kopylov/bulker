package shell

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/model"
	"github.com/pkg/browser"
	"regexp"
	"strings"
)

type RepoPage struct {
	Name         string
	GithubUrl    string
	BitbucketUrl string
}

var (
	PageSources = RepoPage{
		Name:         "sources",
		GithubUrl:    "https://github.com/%v/%v",
		BitbucketUrl: "https://bitbucket.org/%v/%v",
	}
	PageBranches = RepoPage{
		Name:         "branches",
		GithubUrl:    "https://github.com/%v/%v/branches",
		BitbucketUrl: "https://bitbucket.org/%v/%v/branches",
	}
	PageBuilds = RepoPage{
		Name:         "builds",
		GithubUrl:    "https://github.com/%v/%v/actions",
		BitbucketUrl: "https://bitbucket.org/%v/%v/pipelines",
	}
	PagePulls = RepoPage{
		Name:         "pulls",
		GithubUrl:    "https://github.com/%v/%v/pulls",
		BitbucketUrl: "https://bitbucket.org/%v/%v/pull-requests",
	}
)

func (p *RepoPage) String() string {
	return p.Name
}

func (p *RepoPage) Set(v string) error {
	switch v {
	case PageSources.Name:
		*p = PageSources
		return nil
	case PageBranches.Name:
		*p = PageBranches
		return nil
	case PageBuilds.Name:
		*p = PageBuilds
		return nil
	case PagePulls.Name:
		*p = PagePulls
		return nil
	default:
		return fmt.Errorf(
			"must be one of '%s' '%s' '%s' '%s'", PageSources.Name, PageBranches.Name, PageBuilds.Name,
			PagePulls.Name,
		)
	}
}

func (p *RepoPage) Type() string {
	return "RepoPage"
}

type RepoTypeName string

const (
	RepoTypeNameGithubCom    RepoTypeName = "github.com"
	RepoTypeNameBitbucketOrg RepoTypeName = "bitbucket.org"
)

type RepoType interface {
	name() string
	getUrlTemplate(page RepoPage) string
}

type RepoTypeGithubCom struct {
}

func (r *RepoTypeGithubCom) name() string {
	return string(RepoTypeNameGithubCom)
}

func (r *RepoTypeGithubCom) getUrlTemplate(page RepoPage) string {
	return page.GithubUrl
}

type RepoTypeBitbucketOrg struct {
}

func (r *RepoTypeBitbucketOrg) name() string {
	return string(RepoTypeNameBitbucketOrg)
}

func (r *RepoTypeBitbucketOrg) getUrlTemplate(page RepoPage) string {
	return page.BitbucketUrl
}

func OpenPage(repo *model.Repo, page RepoPage) (string, error) {
	repoType, err := getRepoType(repo.Url)
	if err != nil {
		return "", err
	}

	repoTenant, repoId, err := parseUrlDetails(repo.Url)
	if err != nil {
		return "", err
	}

	urlTemplate := repoType.getUrlTemplate(page)

	url := fmt.Sprintf(urlTemplate, repoTenant, repoId)

	err = browser.OpenURL(url)
	if err != nil {
		return "", err
	}

	return url, nil
}

func getRepoType(url string) (RepoType, error) {
	if strings.Contains(url, "github.com") {
		return &RepoTypeGithubCom{}, nil
	} else if strings.Contains(url, "bitbucket.org") {
		return &RepoTypeBitbucketOrg{}, nil
	} else {
		return nil, fmt.Errorf("repository type not supported: %v", url)
	}
}

func parseUrlDetails(url string) (string, string, error) {
	reg, err := regexp.Compile(`.+[:/](.+)/(.+)\.git`)
	if err != nil {
		return "", "", err
	}

	submatch := reg.FindStringSubmatch(url)
	if submatch == nil {
		return "", "", fmt.Errorf("failed to parse url details: %v", url)
	}

	return submatch[1], submatch[2], nil
}

package runner

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"testing"
)

func TestFilter_Matches(t *testing.T) {
	newRepo := func(name string) settings.Repo {
		return settings.Repo{Name: name}
	}
	newRepoWithTags := func(name string, tags []string) settings.Repo {
		return settings.Repo{Name: name, Tags: tags}
	}
	newGroup := func(name string, repos ...string) settings.Group {
		return settings.Group{
			Name:  name,
			Repos: repos,
		}
	}

	tests := []struct {
		name   string
		filter Filter
		repo   settings.Repo
		groups []settings.Group
		want   bool
	}{
		// names
		{
			name: "all", filter: Filter{
				Names: nil,
				Tags:  nil,
			},
			repo: newRepo("qwe"),
			want: true,
		},
		{
			name: "single name", filter: Filter{
				Names: []string{"qwe"},
				Tags:  nil,
			},
			repo: newRepo("qwe"),
			want: true,
		},
		{
			name: "two names", filter: Filter{
				Names: []string{"qwe", "another"},
				Tags:  nil,
			},
			repo: newRepo("qwe"),
			want: true,
		},
		{
			name: "not matching name", filter: Filter{
				Names: []string{"another"},
				Tags:  nil,
			},
			repo: newRepo("qwe"),
			want: false,
		},
		{
			name: "except name", filter: Filter{
				Names: []string{"-qwe"},
				Tags:  nil,
			},
			repo: newRepo("qwe"),
			want: false,
		},
		{
			name: "except another name", filter: Filter{
				Names: []string{"-another"},
				Tags:  nil,
			},
			repo: newRepo("qwe"),
			want: true,
		},
		// tags
		{
			name: "matches tag", filter: Filter{
				Names: nil,
				Tags:  []string{"t1"},
			},
			repo: newRepoWithTags("qwe", []string{"t1", "t2"}),
			want: true,
		},
		{
			name: "matches one of tags", filter: Filter{
				Names: nil,
				Tags:  []string{"t1", "t3"},
			},
			repo: newRepoWithTags("qwe", []string{"t1", "t2"}),
			want: false,
		},
		{
			name: "doesn't match any tag", filter: Filter{
				Names: nil,
				Tags:  []string{"t3", "t4"},
			},
			repo: newRepoWithTags("qwe", []string{"t1", "t2"}),
			want: false,
		},
		{
			name: "except another tag", filter: Filter{
				Names: nil,
				Tags:  []string{"-t3"},
			},
			repo: newRepoWithTags("qwe", []string{"t1", "t2"}),
			want: true,
		},
		{
			name: "except two another tags", filter: Filter{
				Names: nil,
				Tags:  []string{"-t3", "-t4"},
			},
			repo: newRepoWithTags("qwe", []string{"t1", "t2"}),
			want: true,
		},
		{
			name: "except repo tag", filter: Filter{
				Names: nil,
				Tags:  []string{"-t2"},
			},
			repo: newRepoWithTags("qwe", []string{"t1", "t2"}),
			want: false,
		},
		{
			name: "except another tag and except repo one", filter: Filter{
				Names: nil,
				Tags:  []string{"-t2", "-t3"},
			},
			repo: newRepoWithTags("qwe", []string{"t1", "t2"}),
			want: false,
		},
		// both names and tags
		{
			name: "name and tag", filter: Filter{
				Names: []string{"qwe"},
				Tags:  []string{"t1"},
			},
			repo: newRepoWithTags("qwe", []string{"t1", "t2"}),
			want: true,
		},
		{
			name: "name but not tag", filter: Filter{
				Names: []string{"qwe"},
				Tags:  []string{"t3"},
			},
			repo: newRepoWithTags("qwe", []string{"t1", "t2"}),
			want: false,
		},
		{
			name: "not name but tag", filter: Filter{
				Names: []string{"another"},
				Tags:  []string{"t1"},
			},
			repo: newRepoWithTags("qwe", []string{"t1", "t2"}),
			want: false,
		},
		//groups
		{
			name: "group", filter: Filter{
				Names:  []string{},
				Tags:   []string{},
				Groups: []string{"g1"},
			},
			repo:   newRepoWithTags("qwe", []string{}),
			groups: []settings.Group{newGroup("g1", "qwe")},
			want:   true,
		},
		{
			name: "other group", filter: Filter{
				Names:  []string{},
				Tags:   []string{},
				Groups: []string{"g1"},
			},
			repo:   newRepoWithTags("qwe", []string{}),
			groups: []settings.Group{newGroup("g1", "asd")},
			want:   false,
		},
		{
			name: "unknown group", filter: Filter{
				Names:  []string{},
				Tags:   []string{},
				Groups: []string{"g1"},
			},
			repo:   newRepoWithTags("qwe", []string{}),
			groups: []settings.Group{newGroup("g2", "qwe")},
			want:   false,
		},
		{
			name: "not group", filter: Filter{
				Names:  []string{},
				Tags:   []string{},
				Groups: []string{"-g1"},
			},
			repo:   newRepoWithTags("qwe", []string{}),
			groups: []settings.Group{newGroup("g1", "qwe")},
			want:   false,
		},
		{
			name: "not unknown group", filter: Filter{
				Names:  []string{},
				Tags:   []string{},
				Groups: []string{"-g2"},
			},
			repo:   newRepoWithTags("qwe", []string{}),
			groups: []settings.Group{newGroup("g1", "qwe")},
			want:   true,
		},
		{
			name: "two groups", filter: Filter{
				Names:  []string{},
				Tags:   []string{},
				Groups: []string{"g1", "g2"},
			},
			repo:   newRepoWithTags("qwe", []string{}),
			groups: []settings.Group{newGroup("g1", "qwe"), newGroup("g2", "qwe")},
			want:   true,
		},
		{
			name: "belongs only to one group", filter: Filter{
				Names:  []string{},
				Tags:   []string{},
				Groups: []string{"g1", "g2"},
			},
			repo:   newRepoWithTags("qwe", []string{}),
			groups: []settings.Group{newGroup("g1", "qwe"), newGroup("g2", "asd")},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := tt.filter.MatchesRepo(tt.repo, tt.groups); got != tt.want {
					t.Errorf("MatchesRepo() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

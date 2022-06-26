package runner

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"testing"
)

func TestFilter_Matches(t *testing.T) {
	newRepo := func(name string) settings.Repo {
		return settings.Repo{Name: name}
	}

	tests := []struct {
		name   string
		filter Filter
		repo   settings.Repo
		want   bool
	}{
		{
			name: "dummy", filter: Filter{
				Names: nil,
				Tags:  nil,
			},
			repo: newRepo("repo"),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := tt.filter.Matches(tt.repo); got != tt.want {
					t.Errorf("Matches() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

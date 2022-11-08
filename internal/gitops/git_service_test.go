package gitops

import (
	"github.com/mih-kopylov/bulker/internal/shell"
	"reflect"
	"testing"
)

func Test_parseBranches(t *testing.T) {
	tests := []struct {
		name                string
		consoleOutputString string
		want                []Branch
		wantErr             bool
	}{
		{
			name:                "single local branch",
			consoleOutputString: `refs/heads/master`,
			want:                []Branch{{Name: "master"}},
			wantErr:             false,
		},
		{
			name:                "local branch with slashes",
			consoleOutputString: `refs/heads/a/b/c`,
			want:                []Branch{{Name: "a/b/c"}},
			wantErr:             false,
		},
		{
			name:                "single remote branch",
			consoleOutputString: `refs/remotes/origin/master`,
			want:                []Branch{{Name: "master", Remote: "origin"}},
			wantErr:             false,
		},
		{
			name: "remote HEAD is skipped",
			consoleOutputString: `refs/remotes/origin/HEAD
refs/remotes/origin/master`,
			want:    []Branch{{Name: "master", Remote: "origin"}},
			wantErr: false,
		},
		{
			name:                "unsupported branch",
			consoleOutputString: `refs/tags/1.2.3`,
			want:                nil,
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gitService := NewGitService(&shell.NativeShell{})
				got, err := gitService.parseBranches(tt.consoleOutputString)
				if (err != nil) != tt.wantErr {
					t.Errorf("parseBranches() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("parseBranches() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

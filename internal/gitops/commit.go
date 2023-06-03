package gitops

import (
	"fmt"
	"time"
)

type Commit struct {
	Id         string
	AuthorDate time.Time
	Author     CommitUser
	CommitDate time.Time
	Committer  CommitUser
}

type CommitUser struct {
	Name  string
	Email string
}

func (u *CommitUser) String() string {
	return fmt.Sprintf("%v <%v>", u.Name, u.Email)
}

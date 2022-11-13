package gitops

type CloneResult string

const (
	ClonedSuccessfully CloneResult = "cloned"
	ClonedAlready      CloneResult = "already cloned"
	ClonedAgain        CloneResult = "re-cloned"
	CloneError         CloneResult = "error"
)

func (r *CloneResult) String() string {
	return string(*r)
}

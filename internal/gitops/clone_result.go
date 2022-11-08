package gitops

type CloneResult string

const (
	ClonedSuccessfully CloneResult = "Cloned"
	ClonedAlready      CloneResult = "Already cloned"
	ClonedAgain        CloneResult = "Re-cloned"
	CloneError         CloneResult = "Error"
)

func (r *CloneResult) String() string {
	return string(*r)
}

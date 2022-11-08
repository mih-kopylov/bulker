package gitops

type StatusResult string

const (
	StatusClean   StatusResult = "Clean"
	StatusDirty   StatusResult = "Dirty"
	StatusMissing StatusResult = "Missing"
	StatusError   StatusResult = "Error"
)

func (r *StatusResult) String() string {
	return string(*r)
}

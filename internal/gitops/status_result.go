package gitops

type StatusResult string

const (
	StatusClean   StatusResult = "clean"
	StatusDirty   StatusResult = "dirty"
	StatusMissing StatusResult = "missing"
	StatusError   StatusResult = "error"
)

func (r *StatusResult) String() string {
	return string(*r)
}

package gitops

type CreateResult string

const (
	CreateOk    CreateResult = "created"
	CreateError CreateResult = "error"
)

func (r *CreateResult) String() string {
	return string(*r)
}

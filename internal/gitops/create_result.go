package gitops

type CreateResult string

const (
	CreateOk    CreateResult = "Created"
	CreateError CreateResult = "Error"
)

func (r *CreateResult) String() string {
	return string(*r)
}

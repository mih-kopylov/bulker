package gitops

type CheckoutResult string

const (
	CheckoutOk       CheckoutResult = "success"
	CheckoutNotFound CheckoutResult = "not found"
	CheckoutError    CheckoutResult = "error"
)

func (r *CheckoutResult) String() string {
	return string(*r)
}

package gitops

type CheckoutResult string

const (
	CheckoutOk       CheckoutResult = "Success"
	CheckoutNotFound CheckoutResult = "Not Found"
	CheckoutError    CheckoutResult = "Error"
)

func (r *CheckoutResult) String() string {
	return string(*r)
}

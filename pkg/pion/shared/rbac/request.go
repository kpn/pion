package rbac

type Decision string

var (
	PermitDecision Decision = "permit"
	DenyDecision   Decision = "deny"
)

// Request represents the RBAC authorization request
type Request struct {
	UserID string   `json:"userId"` // user identifier, i.e. KPN ruisnaam
	Groups []string `json:"groups"` // available groups that the user belongs to
	Action Action   `json:"action"`
	Target Resource `json:"target"`
}

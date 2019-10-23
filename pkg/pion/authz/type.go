package authz

// ActionType define action values in ACLs
type ActionType string

// DecisionType define decision values of ACLs
type DecisionType string

const (
	Read  ActionType = "Read"
	Write ActionType = "Write"
)

const (
	DecisionPermit DecisionType = "Permit"
	DecisionDeny   DecisionType = "Deny"
)

// Request represents an authorization request
type Request struct {
	Username string     `json:"username"` // user identifier, i.e. KPN ruisnaam
	Groups   []string   `json:"groups"`   // available groups that the user belongs to
	Action   ActionType `json:"action"`   // action to the bucket
	Target   string     `json:"target"`   // name of bucket
	Customer string     `json:"customer"` // customer that owns the bucket
}

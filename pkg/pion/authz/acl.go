package authz

// ACLList is the array of ACLs
type ACLList []ACL

// ACL defines a ACL object used in the bucket
type ACL struct {
	Id       string       `json:"id"`
	Actions  []ActionType `json:"actions"`
	Grantees []Grantee    `json:"grantees"`
}

// ACL evaluation algorithm:
// - Check if the request matches the ACL scope and permission
// - Deny by default
// Evaluate checks the request r against the current ACL
func (acl *ACL) Evaluate(r Request) DecisionType {
	if acl.matchGrantee(r) && acl.matchAction(r) {
		return DecisionPermit
	}

	return DecisionDeny
}

func (acl *ACL) matchGrantee(r Request) bool {
	for _, g := range acl.Grantees {
		if g.Match(r) {
			return true
		}
	}
	return false
}

func (acl *ACL) matchAction(r Request) bool {
	for _, a := range acl.Actions {
		if r.Action == a {
			return true
		}
	}
	return false
}

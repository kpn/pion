package authz

type GranteeType string

const (
	UserType  GranteeType = "User"
	GroupType GranteeType = "Group"
)

type Grantee struct {
	Type  GranteeType `json:"type"`
	Value string      `json:"value"`
}

// Match returns true if the request matches the grantee spec
func (g Grantee) Match(r Request) bool {
	switch g.Type {
	case UserType:
		return g.Value == r.Username
	case GroupType:
		for _, group := range r.Groups {
			if group == g.Value {
				return true
			}
		}
	}
	return false
}

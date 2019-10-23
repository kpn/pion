package rbac

type Action string
type Resource string

type Rule struct {
	Resources []Resource `json:"resources"`
	Actions   []Action   `json:"actions"`
}

// Role represents a role element in the RBAC model
type Role struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Rules       []Rule `json:"rules"`
}

// Match checks if the given <resource, action> matches to the role
func (role Role) Match(a Action, r Resource) bool {
	for _, rule := range role.Rules {
		if matchResource(rule.Resources, r) && matchAction(rule.Actions, a) {
			return true
		}
	}

	return false
}

func matchResource(resources []Resource, target Resource) bool {
	for _, r := range resources {
		if r == target {
			return true
		}
	}
	return false
}

func matchAction(actions []Action, target Action) bool {
	for _, a := range actions {
		if a == target {
			return true
		}
	}
	return false
}

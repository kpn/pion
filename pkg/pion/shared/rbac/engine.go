package rbac

import "github.com/golang/glog"

// Engine defines the API for RBAC engine
type Engine interface {
	// Evaluate returns the decision of the engine against the given request
	Evaluate(r Request) Decision
}

type rbacEngine struct {
	roleBindings []RoleBinding
	rolesMap     map[string]Role // mapping from role name to role object

	mappingRolesFromUser  map[string][]string // map from userId to rolesMap
	mappingRolesFromGroup map[string][]string // map from group to rolesMap
}

// NewEngine creates a new RBAC evaluation engine with defined roles and role-bindings
func NewEngine(roles []Role, roleBindings []RoleBinding) Engine {
	engine := &rbacEngine{
		roleBindings: roleBindings,
	}

	engine.rolesMap = make(map[string]Role)
	for _, r := range roles {
		engine.rolesMap[r.Name] = r
	}
	engine.initRoleBindings()
	return engine
}

func (e *rbacEngine) initRoleBindings() {
	e.mappingRolesFromUser = make(map[string][]string)
	e.mappingRolesFromGroup = make(map[string][]string)

	for _, rb := range e.roleBindings {
		for _, subj := range rb.Subjects {
			if subj.Type == UserType {
				e.mappingRolesFromUser[subj.Value] = append(e.mappingRolesFromUser[subj.Value], rb.RoleRef)
			}
			if subj.Type == GroupType {
				e.mappingRolesFromGroup[subj.Value] = append(e.mappingRolesFromGroup[subj.Value], rb.RoleRef)
			}
		}
	}
}

func (e rbacEngine) Evaluate(r Request) Decision {
	validRoleNames := e.findValidRoleNames(r)
	glog.V(2).Infof("Valid roles matching the request: %v", validRoleNames)

	for _, rName := range validRoleNames {
		glog.V(3).Infof("Checking role '%s'", rName)
		role, ok := e.rolesMap[rName]
		if ok && role.Match(r.Action, r.Target) {
			glog.V(3).Infof("Request matched role '%s' - permitted", rName)
			return PermitDecision
		}
	}
	glog.V(3).Info("No matching role, denied")
	return DenyDecision
}

func (e rbacEngine) findValidRoleNames(r Request) []string {
	validRoleNames := e.mappingRolesFromUser[r.UserID]
	for _, g := range r.Groups {
		validRoleNames = append(validRoleNames, e.mappingRolesFromGroup[g]...)
	}
	return unique(validRoleNames)
}

func unique(slice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range slice {
		if _, existed := keys[entry]; !existed {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

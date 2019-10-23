package rbac

import (
	"errors"
)

// TODO manage roles in etcd store (TIPC-958)
var DefaultRoles = []Role{
	ReadOnlyRole,
	UserRole,
	EditorRole,
	AdminRole,
}

const (
	BucketResource      Resource = "oss:bucket"
	BucketAclResource   Resource = "oss:bucket-acl"
	ObjectResource      Resource = "oss:object"
	RoleResource        Resource = "oss:role"
	RoleBindingResource Resource = "oss:role-binding"
	CustomerResource    Resource = "oss:customer"

	UnknownResource Resource = "oss:unknown"
)

const (
	Create Action = "create"
	Update Action = "update"
	List   Action = "list"
	Get    Action = "get"
	Delete Action = "delete"

	Publish   Action = "publish"   // Action specific for bucket only, publish a path in the buckets
	Unpublish Action = "unpublish" // un-publish a path in the bucket
)

var (
	ReadOnlyRole = Role{
		Name:        "readonly",
		DisplayName: "ReadOnly",
		Rules: []Rule{
			{
				Resources: []Resource{BucketResource, ObjectResource},
				Actions:   []Action{List, Get},
			},
		},
	}

	UserRole = Role{
		Name:        "user",
		DisplayName: "User",
		Rules: []Rule{
			{
				Resources: []Resource{BucketResource},
				Actions:   []Action{List, Get},
			},
			{
				Resources: []Resource{ObjectResource},
				Actions:   []Action{List, Get, Create, Delete, Update},
			},
		},
	}

	EditorRole = Role{
		Name:        "editor",
		DisplayName: "Editor",
		Rules: []Rule{
			{
				Resources: []Resource{BucketResource},
				Actions:   []Action{List, Get, Create, Delete, Update, Publish, Unpublish},
			},
			{
				Resources: []Resource{ObjectResource},
				Actions:   []Action{List, Get, Create, Delete, Update},
			},
		},
	}

	AdminRole = Role{
		Name:        "admin",
		DisplayName: "Admin",
		Rules: []Rule{
			{
				Resources: []Resource{BucketResource},
				Actions:   []Action{List, Get, Create, Delete, Update, Publish, Unpublish},
			},
			{
				Resources: []Resource{ObjectResource},
				Actions:   []Action{List, Get, Create, Delete, Update},
			},
			{
				Resources: []Resource{RoleResource}, // for now Roles are hard-coded
				Actions:   []Action{List, Get},
			},
			{
				Resources: []Resource{BucketAclResource},
				Actions:   []Action{List, Get, Create, Delete, Update},
			},
			{
				Resources: []Resource{RoleBindingResource},
				Actions:   []Action{List, Get, Update, Create, Delete},
			},
		},
	}
)

func Str2Action(str string) (Action, error) {
	var actionMap = map[string]Action{
		string(Create):    Create,
		string(Update):    Update,
		string(List):      List,
		string(Get):       Get,
		string(Delete):    Delete,
		string(Publish):   Publish,
		string(Unpublish): Unpublish,
	}
	action := actionMap[str]
	if action == "" {
		return "", errors.New("unknown action")
	}
	return action, nil
}

func Str2Resource(str string) (Resource, error) {
	var resourceMap = map[string]Resource{
		string(BucketResource):      BucketResource,
		string(BucketAclResource):   BucketAclResource,
		string(ObjectResource):      ObjectResource,
		string(RoleResource):        RoleResource,
		string(RoleBindingResource): RoleBindingResource,
		string(CustomerResource):    CustomerResource,
	}
	res := resourceMap[str]
	if res == "" {
		return "", errors.New("unknown resource")
	}
	return res, nil
}

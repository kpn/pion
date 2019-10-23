package rbac_test

import (
	"testing"

	"github.com/kpn/pion/pkg/pion/shared/rbac"
)

func TestAdminRole(t *testing.T) {
	roleBinding := rbac.RoleBinding{
		Name:    "default-admin-rb",
		RoleRef: "admin",
		Subjects: []rbac.Subject{
			{Type: rbac.GroupType, Value: "dig_infraplatform"},
			{Type: rbac.UserType, Value: "ngo500"},
		},
	}

	testCases := []struct {
		Request  rbac.Request
		Decision rbac.Decision
	}{
		{
			rbac.Request{
				UserID: "foo",
				Groups: []string{"dig_infraplatform"},
				Action: rbac.Update,
				Target: rbac.RoleBindingResource,
			},
			rbac.PermitDecision,
		},
		{
			rbac.Request{
				UserID: "foo",
				Groups: []string{"dig_infraplatform"},
				Action: rbac.Delete,
				Target: rbac.RoleBindingResource,
			},
			rbac.PermitDecision,
		},
		{
			rbac.Request{
				UserID: "ngo500",
				Groups: []string{"dig_foo"},
				Action: rbac.Create,
				Target: rbac.RoleBindingResource,
			},
			rbac.PermitDecision,
		},
	}

	engine := rbac.NewEngine(
		[]rbac.Role{rbac.UserRole,
			rbac.EditorRole,
			rbac.AdminRole},
		[]rbac.RoleBinding{roleBinding})

	for _, c := range testCases {
		actual := engine.Evaluate(c.Request)
		expected := c.Decision
		if expected != actual {
			t.Errorf("Evaluation failed, expected '%v', actual '%v':\nRequest: '%+v'", expected, actual, c.Request)
		}
	}
}

func TestEditorRole(t *testing.T) {
	roleBinding := rbac.RoleBinding{
		Name:    "default-poweruser-rb",
		RoleRef: "editor",
		Subjects: []rbac.Subject{
			{Type: rbac.GroupType, Value: "dig_foo"},
			{Type: rbac.UserType, Value: "ngo500"},
		},
	}

	testCases := []struct {
		Request  rbac.Request
		Decision rbac.Decision
	}{
		{
			rbac.Request{
				UserID: "bah",
				Groups: []string{"dig_foo"},
				Action: rbac.Publish,
				Target: rbac.BucketResource,
			},
			rbac.PermitDecision,
		},
		{
			rbac.Request{
				UserID: "bah",
				Groups: []string{"dig_foo"},
				Action: rbac.Unpublish,
				Target: rbac.BucketResource,
			},
			rbac.PermitDecision,
		},
		{
			rbac.Request{
				UserID: "ngo500",
				Groups: []string{"dig_random"},
				Action: rbac.Create,
				Target: rbac.BucketResource,
			},
			rbac.PermitDecision,
		},
		{
			rbac.Request{
				UserID: "ngo500",
				Groups: []string{"dig_random"},
				Action: rbac.Delete,
				Target: rbac.BucketResource,
			},
			rbac.PermitDecision,
		},
		{
			rbac.Request{
				UserID: "ngo500",
				Groups: []string{"dig_random"},
				Action: rbac.Get,
				Target: rbac.RoleBindingResource,
			},
			rbac.DenyDecision,
		},
	}

	engine := rbac.NewEngine(
		[]rbac.Role{rbac.UserRole,
			rbac.EditorRole,
			rbac.AdminRole},
		[]rbac.RoleBinding{roleBinding})

	for _, c := range testCases {
		actual := engine.Evaluate(c.Request)
		expected := c.Decision
		if expected != actual {
			t.Errorf("Evaluation failed, expected '%v', actual '%v':\nRequest: '%+v'", expected, actual, c.Request)
		}
	}
}

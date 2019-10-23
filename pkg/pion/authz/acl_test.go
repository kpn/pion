package authz_test

import (
	"testing"

	"github.com/kpn/pion/pkg/pion/authz"
)

func TestACLEval(t *testing.T) {
	testCases := []struct {
		Acl       authz.ACL
		Requests  []authz.Request
		Decisions []authz.DecisionType
	}{
		{
			Acl: authz.ACL{
				Id:      "1",
				Actions: []authz.ActionType{authz.Read, authz.Write},
				Grantees: []authz.Grantee{
					{
						authz.UserType,
						"alice",
					},
					{
						authz.GroupType,
						"alpha",
					},
					{
						authz.GroupType,
						"beta",
					},
				},
			},
			Requests: []authz.Request{
				{
					Username: "cngo",
					Groups:   []string{"alpha", "gamma"},
					Action:   authz.Read,
					Target:   "foo-bucket",
				},
				{
					Username: "cngo",
					Groups:   []string{"beta", "omega"},
					Action:   authz.Write,
					Target:   "foo-bucket",
				},
			},
			Decisions: []authz.DecisionType{
				authz.DecisionPermit,
				authz.DecisionPermit,
			},
		},
		{
			Acl: authz.ACL{
				Id:      "2",
				Actions: []authz.ActionType{authz.Read},
				Grantees: []authz.Grantee{
					{
						authz.UserType,
						"alice",
					},
					{
						authz.GroupType,
						"alpha",
					},
					{
						authz.GroupType,
						"beta",
					},
				},
			},
			Requests: []authz.Request{
				{
					Username: "cngo",
					Groups:   []string{"alpha", "gamma"},
					Action:   authz.Write,
					Target:   "foo-bucket",
				},
				{
					Username: "cngo",
					Groups:   []string{"beta", "omega"},
					Action:   authz.Read,
					Target:   "foo-bucket",
				},
			},
			Decisions: []authz.DecisionType{
				authz.DecisionDeny,
				authz.DecisionPermit,
			},
		},
	}

	for ci, c := range testCases {
		for ri, r := range c.Requests {
			if expected, actual := c.Decisions[ri], c.Acl.Evaluate(r); expected != actual {
				t.Errorf("test-case a%d-r%d: evaluation failed, expected '%v', actual '%v'\nAcl: '%+v'\nRequest: '%+v' ", ci, ri, expected, actual, c.Acl, r)
			}

		}
	}
}

package authz_test

import (
	"testing"

	"github.com/kpn/pion/pkg/pion/authz"
)

func TestMatchGrantee(t *testing.T) {
	testCases := []struct {
		Grantee  authz.Grantee
		Request  authz.Request
		Decision bool
	}{
		{
			authz.Grantee{
				Type:  authz.UserType,
				Value: "cngo",
			},
			authz.Request{
				Username: "cngo",
				Groups:   []string{"alpha", "beta"},
			},
			true,
		},
		{
			authz.Grantee{
				Type:  authz.UserType,
				Value: "cngo",
			},
			authz.Request{
				Username: "foo",
				Groups:   []string{"alpha", "beta"},
			},
			false,
		},
		{
			authz.Grantee{
				Type:  authz.GroupType,
				Value: "alpha",
			},
			authz.Request{
				Username: "foo",
				Groups:   []string{"alpha", "beta"},
			},
			true,
		},
		{
			authz.Grantee{
				Type:  authz.GroupType,
				Value: "gamma",
			},
			authz.Request{
				Username: "foo",
				Groups:   []string{"alpha", "beta"},
			},
			false,
		},
		{
			authz.Grantee{
				Type:  authz.GroupType,
				Value: "foo",
			},
			authz.Request{
				Username: "foo",
				Groups:   []string{"alpha", "beta"},
			},
			false,
		},
	}

	for i, c := range testCases {
		if expected, actual := c.Decision, c.Grantee.Match(c.Request); actual != expected {
			t.Errorf("test-case %d: matching failed, expected '%v', actual '%v'\nRequest: '%+v'\nGrantee: '%+v' ", i, expected, actual, c.Request, c.Grantee)
		}
	}
}

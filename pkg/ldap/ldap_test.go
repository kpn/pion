package ldap_test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/kpn/pion/pkg/ldap"
	"github.com/kpn/pion/pkg/pion"
	"github.com/stretchr/testify/assert"
)

func init() {
	flag.Set("alsologtostderr", fmt.Sprintf("%t", true))
	var logLevel string
	flag.StringVar(&logLevel, "logLevel", "2", "test")
	flag.Lookup("v").Value.Set(logLevel)
}

func TestNewLDAPClient(t *testing.T) {
	// TODO exec slapd with sample database prior running this test
	t.SkipNow()
	config := pion.LDAPSetting{
		Host:         "localhost",
		Port:         3389,
		Type:         "ldap",
		BindDN:       "cn=admin,dc=example,dc=org",
		BindPassword: "QoyzmJVZv41BsANrTldJjl23EOZGIKyy",
		UserBase:     "ou=People,dc=example,dc=org",
		GroupBase:    "ou=Groups,dc=example,dc=org",
		UserClass:    "inetOrgPerson",
		GroupClass:   "groupofnames",
	}

	client, err := ldap.NewLDAPClient(config)
	assert.NoError(t, err)

	err = client.Login("billy", "billy")
	assert.NoError(t, err)

	user, err := client.GetUser("billy")
	assert.NoError(t, err)
	fmt.Println(user)

	groups, err := client.GetUserGroups("billy")
	assert.NoError(t, err)
	fmt.Println(groups)
	assert.Len(t, groups, 2)
	assert.Contains(t, groups, "group_customer1")
	assert.Contains(t, groups, "group_customer2")
}

package ldap

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// GetUser queries LDAP and return the user profile
func (c LDAPClient) GetUser(username string) (*UserInfo, error) {
	query := fmt.Sprintf("(&(uid=%s)(objectClass=%s))", username, c.setting.UserClass)
	attrs := []string{"uid", "sn", "cn", "mail"}

	result, err := c.QueryLDAP(query, attrs, c.setting.UserBase)
	if err != nil {
		return nil, err
	}

	if len(result.Entries) == 0 {
		glog.V(2).Infof("User '%s' not found", username)
		return nil, nil
	}

	if len(result.Entries) > 1 {
		glog.Warningf("Multiple user entries for '%s'. Something wrong?", username)
		return nil, errors.New("multiple user entries found")
	}

	entry := result.Entries[0]
	return &UserInfo{
		Username:  entry.GetAttributeValue("uid"),
		FirstName: entry.GetAttributeValue("cn"),
		LastName:  entry.GetAttributeValue("sn"),
		Mail:      entry.GetAttributeValue("mail"),
		// Title:       entry.GetAttributeValue("title"),s
		// DisplayName: entry.GetAttributeValue("displayName"),
	}, nil
}

// GetUserGroups returns list of groups that user is the member of
func (c LDAPClient) GetUserGroups(username string) (groups []string, err error) {
	query := fmt.Sprintf("(member=cn=%s,%s)", username, c.setting.UserBase)
	glog.V(1).Info("Find groups of user with query:", query)
	attrs := []string{"cn"}

	result, err := c.QueryLDAP(query, attrs, c.setting.GroupBase)
	if err != nil {
		return nil, err
	}

	for _, entry := range result.Entries {
		groups = append(groups, entry.GetAttributeValue("cn"))
	}
	return groups, nil
}

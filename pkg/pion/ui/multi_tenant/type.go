package multi_tenant

import "github.com/kpn/pion/pkg/ldap"

type UserInfo struct {
	ldap.UserInfo
	Customer   string   `json:"customer"`   // customer account that the user belongs to
	UserGroups []string `json:"userGroups"` // user-group of the user that belongs in the customer account
}

// This is the list of user attributes used in multi-tenant model
const (
	// UserAttributeCustomer attribute refers to the customer of the user
	UserAttributeCustomer = "customer"

	// UserAttributeGroups attribute refers to the groups value of the user
	UserAttributeGroups = "groups"
)

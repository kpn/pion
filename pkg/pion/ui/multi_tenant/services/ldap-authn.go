package services

import (
	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/ldap"
	"github.com/kpn/pion/pkg/pion-clients/manager-client"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/kpn/pion/pkg/pion/ui/multi_tenant"
	"github.com/kpn/pion/pkg/pion/ui/util"
	"github.com/pkg/errors"
)

var (
	ErrCustomerNotFound  = errors.New("customer not found")
	ErrUserNotInCustomer = errors.New("user not in customer")
)

type ldapAuthnService struct {
}

func NewLDAPAuthnService() AuthnService {
	return &ldapAuthnService{}
}

func (svc ldapAuthnService) Login(customerName string, username, password string) (*multi_tenant.UserInfo, error) {
	cc := manager_client.NewCustomerClient(manager_client.ManagerServiceURL)
	customer, err := cc.Get(customerName)
	if err != nil {
		glog.Warningf("Failed to get customer: %v", err)
		return nil, ErrCustomerNotFound
	}

	ldapUser, err := ldap.Login(username, password)
	if err != nil {
		return nil, err
	}

	groupsInCustomer, err := groupsInCustomer(ldapUser.Username, customer)
	if err != nil {
		return nil, err
	}
	if len(groupsInCustomer) == 0 && !customer.ContainsUser(ldapUser.Username) {
		return nil, ErrUserNotInCustomer
	}

	return &multi_tenant.UserInfo{
		UserInfo:   *ldapUser,
		Customer:   customer.Name,
		UserGroups: groupsInCustomer,
	}, nil
}

// groupsInCustomer checks if the user belongs to the given customer. It returns user's group(s) that in the customer account.
// otherwise, it returns empty group list
func groupsInCustomer(username string, customer *model.Customer) (groups []string, err error) {
	userGroups, err := util.GetUserGroupsFromLDAP(username)
	if err != nil {
		glog.Errorf("Cannot get user group from LDAP: %v", err)
		return nil, err
	}

	// trivial intersection loop, O(n^2)
	var groupsInCustomer []string
	for _, ug := range userGroups {
		for _, cg := range customer.Groups {
			if ug == cg {
				groupsInCustomer = append(groupsInCustomer, ug)
			}
		}
	}
	return groupsInCustomer, nil
}

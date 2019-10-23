package ldap

import (
	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion"
)

func Login(username, password string) (*UserInfo, error) {
	client, err := NewLDAPClient(pion.AppConfig().LDAP)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	err = client.Login(username, password)
	if err != nil {
		glog.V(2).Infof("Failed to bind LDAP for user '%s':%v", username, err)
		return nil, err
	}

	glog.V(2).Infof("Authenticated user '%s' against LDAP", username)

	userInfo, err := client.GetUser(username)
	if err != nil {
		glog.Errorf("Failed to query user information: %v", err)
		return nil, err
	}

	return userInfo, nil
}

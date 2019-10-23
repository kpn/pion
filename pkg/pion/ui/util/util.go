package util

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/ldap"
	"github.com/kpn/pion/pkg/pion"
	"github.com/labstack/echo"
)

func RedirectTo(c echo.Context, page string) error {
	var url string
	if page == "" || page == "/" {
		url = fmt.Sprintf("%s://%s", c.Scheme(), c.Request().Host)
	} else {
		url = fmt.Sprintf("%s://%s/%s", c.Scheme(), c.Request().Host, page)
	}

	return c.Redirect(http.StatusMovedPermanently, url)
}

func GetUserGroupsFromLDAP(username string) ([]string, error) {
	ldapClient, err := ldap.NewAndBindLDAPClient(pion.AppConfig().LDAP)
	if err != nil {
		glog.Errorf("Failed to connect to LDAP server: %v", err)
		return nil, errors.New("failed to connect to LDAP server")
	}
	defer ldapClient.Close()

	// TODO cache userGroups for limited time. Access token will not contains group, but when verifying token,
	// user groups are fetched from the cache (not from LDAP directly)
	return ldapClient.GetUserGroups(username)
}

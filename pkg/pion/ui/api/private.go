package api

import (
	"errors"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion-clients/manager-client"
	"github.com/kpn/pion/pkg/pion/ui/session"
	"github.com/labstack/echo"
)

func createManagerClient(c echo.Context) (manager_client.ManagerClient, error) {
	userInfo, err := session.GetUserInfo(c)
	if err != nil {
		glog.Errorf("Failed to get userInfo from session: %v", err)
		return nil, errors.New("session does not have userInfo object")
	}

	mc := manager_client.New(manager_client.ManagerServiceURL, userInfo.Customer, userInfo.Username, userInfo.UserGroups)
	return mc, nil
}

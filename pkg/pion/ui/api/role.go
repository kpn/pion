package api

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion-clients/manager-client"
	"github.com/kpn/pion/pkg/pion/shared/rbac"
	"github.com/labstack/echo"
)

func ListRoles(c echo.Context) error {
	mc, err := createManagerClient(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	roles, err := mc.ListRoles()
	switch err {
	case manager_client.ErrForbidden:
		return echo.ErrForbidden
	case manager_client.ErrUnauthorized:
		return echo.ErrUnauthorized
	case manager_client.ErrCustomerNotfound:
		return echo.NewHTTPError(http.StatusBadRequest, "customer not found")
	}
	if err != nil {
		glog.Errorf("Failed to list roles: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list roles")
	}
	glog.V(2).Infof("Listing roles: %v", roles)
	return c.JSON(http.StatusOK, roles)
}

func GetRole(c echo.Context) error {
	roleName := c.Param("name")
	if roleName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing name param")
	}

	mc, err := createManagerClient(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	roles, err := mc.ListRoles()
	switch err {
	case manager_client.ErrForbidden:
		return echo.ErrForbidden
	case manager_client.ErrUnauthorized:
		return echo.ErrUnauthorized
	case manager_client.ErrCustomerNotfound:
		return echo.NewHTTPError(http.StatusBadRequest, "customer not found")
	}
	if err != nil {
		glog.Errorf("Failed to get role: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list roles")
	}
	var role *rbac.Role = nil
	for _, r := range roles {
		if r.Name == roleName {
			role = &r
			break
		}
	}
	if role == nil {
		glog.Warningf("Role '%s' not found", roleName)
		return echo.NewHTTPError(http.StatusNotFound, "Role not found")
	}
	glog.V(2).Infof("Get role: %v", role)
	return c.JSON(http.StatusOK, role)
}

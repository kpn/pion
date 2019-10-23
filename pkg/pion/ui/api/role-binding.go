package api

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion-clients/manager-client"
	"github.com/kpn/pion/pkg/pion/shared/rbac"
	"github.com/kpn/pion/pkg/sts/handlers"
	"github.com/labstack/echo"
)

// ListRoleBindings handles API requests to list role-binding objects
func ListRoleBindings(c echo.Context) error {
	mc, err := createManagerClient(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	rbs, err := mc.ListRoleBindings()
	switch err {
	case manager_client.ErrForbidden:
		return echo.ErrForbidden
	case manager_client.ErrUnauthorized:
		return echo.ErrUnauthorized
	case manager_client.ErrCustomerNotfound:
		return echo.NewHTTPError(http.StatusBadRequest, "customer not found")
	}
	if err != nil {
		glog.Errorf("Failed to list role-bindings: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list role-bindings")
	}
	glog.V(2).Infof("Listing role-bindings: %v", rbs)
	return c.JSON(http.StatusOK, rbs)
}

// CreateRoleBinding handles API requests to create role-binding objects
func CreateRoleBinding(c echo.Context) error {
	var valueRegex = regexp.MustCompile(`^\w[\w\-\_]{3,}$`)
	var payload rbac.RoleBinding
	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "parsing payload failed")
	}
	// validate inputs
	if !valueRegex.MatchString(payload.Name) {
		return c.JSON(http.StatusBadRequest, "Invalid name")
	}
	for _, s := range payload.Subjects {
		if !valueRegex.MatchString(s.Value) {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid subject value '%s'", s.Value))
		}
	}

	mc, err := createManagerClient(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	newRB, err := mc.CreateRoleBinding(payload)
	if err != nil {
		glog.Errorf("Create role-binding failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create new role-binding")
	}
	return c.JSON(http.StatusOK, newRB)
}

// DeleteRoleBinding handles API requests to delete role-binding objects
func DeleteRoleBinding(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing 'name' parameter")
	}
	mc, err := createManagerClient(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err = mc.DeleteRoleBinding(name)
	if err == manager_client.ErrRoleBindingNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return handlers.Response(c, http.StatusOK, "Deleted role-binding successfully")
}

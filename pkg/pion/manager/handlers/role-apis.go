package handlers

import (
	"net/http"

	"github.com/kpn/pion/pkg/pion/shared/rbac"
	"github.com/labstack/echo"
)

// ListRoles returns predefined roles
func ListRoles(c echo.Context) error {
	// TODO Manage roles in etcd store rather than hard-coded (TIPC-958)
	return c.JSON(http.StatusOK, rbac.DefaultRoles)
}

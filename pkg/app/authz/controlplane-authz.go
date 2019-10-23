package authz

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/authz/handlers"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/rbac"
	"github.com/labstack/echo"
)

// ControlPlaneAuthorize implements access control for control plane requests, i.e. multi-tenant RBAC
func (app *App) ControlPlaneAuthorize(c echo.Context) error {
	httpRequest := c.Request()
	customerName := httpRequest.Header.Get(shared.CustomerKey)
	if customerName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "customer header not found")
	}

	rh, err := handlers.NewRBACHandler(shared.DefaultEtcdAddress, customerName)
	if err != nil {
		glog.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if rh.AuthorizeHTTP(httpRequest) == rbac.PermitDecision {
		glog.V(2).Infof("Authorization permitted: %+v", httpRequest)
		return c.NoContent(http.StatusOK)
	}
	glog.V(2).Infof("Authorization denied: %v", httpRequest)
	return c.NoContent(http.StatusForbidden)
}

package middleware

import (
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/kpn/pion/pkg/pion/shared/rbac"
	"github.com/labstack/echo"
	"gopkg.in/resty.v1"
)

var authzService = os.Getenv("AUTHZ_SERVICE_URL")

// AuthorizeHTTP middleware invokes the MT-RBAC authz endpoint in the authz-service
func Authorize(resource rbac.Resource, action rbac.Action, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		glog.V(2).Infof("Authorizing action '%s' on '%s'", action, resource)

		httpRequest := c.Request()
		userId := httpRequest.Header.Get(shared.UserIdKey)
		userGroups := httpRequest.Header.Get(shared.UserGroupKey)
		customer, ok := c.Get(shared.CustomerKey).(*model.Customer)
		// userGroups can be empty in case of individual users
		if userId == "" || !ok || customer == nil {
			glog.Error("Unauthenticated request, missing either userId, userGroups or customer attribute")
			return echo.NewHTTPError(http.StatusBadRequest, "missing authenticated data")
		}
		glog.V(3).Infof("userId=%s,groups=%s,action=%s,resource=%s,customer=%s",
			userId, userGroups, action, resource, customer.Name)

		authzEndpoint := authzService + "/mt-rbac-authorize"
		glog.V(3).Infof("Sending request to '%s'", authzEndpoint)

		resp, err := resty.R().
			SetHeader(shared.UserIdKey, userId).
			SetHeader(shared.UserGroupKey, userGroups).
			SetHeader(shared.ActionKey, string(action)).
			SetHeader(shared.ResourceKey, string(resource)).
			SetHeader(shared.CustomerKey, customer.Name).
			Get(authzEndpoint)

		if err != nil {
			glog.Errorf("Cannot call authorization service: %v", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		if resp.StatusCode() != http.StatusOK {
			glog.Warningf("Unauthorized decision from authz-service: %d", resp.StatusCode())
			return c.NoContent(http.StatusForbidden)
		}
		glog.V(3).Info("Granted permission")
		return next(c)
	}
}

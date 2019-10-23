package middleware

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/labstack/echo"
)

// ValidateUserInCustomer checks if the user referring by UserID or Groups attributes in request headers belongs to the
// customer in the request context. If affirmative, it adds userId and groups attribute to request context
func ValidateUserInCustomer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		customer, ok := c.Get(shared.CustomerKey).(*model.Customer)
		if !ok || customer == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "no customer found in the request context")
		}
		userId, userGroups := getUserAttributes(c.Request())
		if userId == "" || len(userGroups) == 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, "userId and/or userGroups attributes not found")
		}

		glog.V(3).Infof("Checking if any of user request groups '%v' exists in customer's groups '%v'", userGroups, customer.Groups)
		if !customer.HasAnyGroup(userGroups) && !customer.ContainsUser(userId) {
			glog.Warningf("User '%s' with groups '%v' does not belong to customer '%s'", userId, userGroups, customer.Name)
			return echo.NewHTTPError(http.StatusForbidden, "user does not belong to the customer")
		}

		// set authenticated attributes to the request context
		c.Set(shared.UserIdKey, userId)
		c.Set(shared.UserGroupKey, userGroups)
		return next(c)
	}
}

package api

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/ui/multi_tenant"
	"github.com/kpn/pion/pkg/pion/ui/multi_tenant/services"
	"github.com/kpn/pion/pkg/pion/ui/session"
	"github.com/kpn/pion/pkg/sts/handlers"
	"github.com/labstack/echo"
)

// Login handles form-based POST requests to login
func Login(c echo.Context) error {
	type Payload struct {
		Customer string `json:"customer"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var p Payload
	err := c.Bind(&p)
	if err != nil {
		glog.Warningf("Failed to parse login payload: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to parse payload")
	}

	if p.Username == "" || p.Password == "" {
		glog.Warningf("Username or password is empty")
		return handlers.Response(c, http.StatusBadRequest, "Username or password is empty")
	}
	glog.Infof("User '%s' is logging in", p.Username)

	var userInfo *multi_tenant.UserInfo
	authnsvc := services.NewLDAPAuthnService()
	userInfo, err = authnsvc.Login(p.Customer, p.Username, p.Password)
	switch err {
	case services.ErrCustomerNotFound:
		fallthrough
	case services.ErrUserNotInCustomer:
		glog.Warning(err.Error())
		return c.NoContent(http.StatusUnauthorized)
	}
	if err != nil {
		glog.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	glog.Infof("User '%s' logged in", userInfo.Username)

	err = session.Create(c, userInfo)
	if err != nil {
		return err
	}
	glog.V(2).Infof("Created a new session for user '%s'", p.Username)
	return c.JSON(http.StatusOK, userInfo)
}

func LogOut(c echo.Context) error {
	err := session.Clear(c)
	if err != nil {
		glog.Errorf("Logout failed: %v", err)
		return handlers.Response(c, http.StatusInternalServerError, "cannot clear session")
	}
	return c.NoContent(http.StatusOK)
}

package session

import (
	"errors"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/kpn/pion/pkg/pion/ui/multi_tenant"
	"github.com/kpn/pion/pkg/pion/ui/util"
	"github.com/labstack/echo"
)

const CookieName = "pion-session"

var (
	authnKey      = securecookie.GenerateRandomKey(32)
	encryptionKey = securecookie.GenerateRandomKey(32)

	// TODO Use a KV-backed session store
	CookieStore = sessions.NewCookieStore(authnKey, encryptionKey)
)

func Get(c echo.Context) (*sessions.Session, error) {
	return CookieStore.Get(c.Request(), CookieName)
}

// Create generates a new authenticated session for given username
func Create(c echo.Context, userInfo *multi_tenant.UserInfo) error {
	s, _ := Get(c)

	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   true,
		// SameSite: http.SameSiteLaxMode,
	}

	s.Values["authenticated"] = true
	s.Values["username"] = userInfo.Username
	s.Values["userinfo"] = *userInfo
	s.Values["customer"] = userInfo.Customer

	err := s.Save(c.Request(), c.Response())
	if err != nil {
		glog.Errorf("Failed to save session: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return nil
}

// Clear the current session
func Clear(c echo.Context) error {
	s, _ := Get(c)
	username := s.Values["username"]
	delete(s.Values, "authenticated")
	delete(s.Values, "username")
	delete(s.Values, "userinfo")

	// delete client cookie
	s.Options.MaxAge = -1

	err := s.Save(c.Request(), c.Response())
	if err != nil {
		glog.Errorf("Failed to save session: %v", err)
		return err
	}
	glog.V(2).Infof("user '%s' logged out", username)
	return nil
}

// Validate is the middleware to check if attached cookie in the request is valid
func Validate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s, _ := CookieStore.Get(c.Request(), CookieName)

		if authn, ok := s.Values["authenticated"].(bool); !ok || !authn {
			glog.V(2).Info("Unauthenticated request")
			return util.RedirectTo(c, "login")
		}

		if s.Values["username"] == "" {
			glog.Error("authenticated session, but user not found")
			return c.NoContent(http.StatusInternalServerError)
		}

		// request contains authenticated cookie, call next
		return next(c)
	}
}

// GetUserInfo extracts userInfo object from the server session
func GetUserInfo(c echo.Context) (*multi_tenant.UserInfo, error) {
	s, _ := Get(c)
	ui, ok := s.Values["userinfo"].(multi_tenant.UserInfo)
	if !ok {

		return nil, errors.New("invalid type casting")
	}
	return &ui, nil
}

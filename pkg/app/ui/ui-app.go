package ui

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/gorilla/context"
	_ "github.com/gorilla/sessions"
	"github.com/kpn/pion/pkg/pion/ui/api"
	"github.com/kpn/pion/pkg/pion/ui/session"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type App struct {
}

func NewApp() App {
	return App{}
}

func (app App) Start() {
	glog.Info("Pion Object Storage Service UI")

	e := echo.New()

	e.Use(clearHandler)
	e.HideBanner = true

	distDir := filepath.Join(os.Getenv("DIST_DIR"))
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   distDir,
		Index:  "index.html",
		HTML5:  true,
		Browse: false,
	}))

	e.POST("/public/login", api.Login)

	// authenticated APIs
	restrictedAPIs := e.Group("/restricted", session.Validate)
	restrictedAPIs.GET("/access_keys", api.ListAccessKeys)
	restrictedAPIs.DELETE("/access_keys/:id", api.DeleteAccessKey)
	restrictedAPIs.POST("/access_keys", api.CreateAccessKey)

	restrictedAPIs.GET("/buckets", api.ListBuckets)
	restrictedAPIs.PUT("/buckets/:id/acl", api.UpdateBucketACLs)
	restrictedAPIs.POST("/buckets", api.AddBucket)

	restrictedAPIs.GET("/roles", api.ListRoles)
	restrictedAPIs.GET("/roles/:name", api.GetRole)

	restrictedAPIs.GET("/role_bindings", api.ListRoleBindings)
	restrictedAPIs.POST("/role_bindings", api.CreateRoleBinding)
	restrictedAPIs.DELETE("/role_bindings/:name", api.DeleteRoleBinding)

	restrictedAPIs.POST("/logout", api.LogOut)
	e.GET("/health", checkHealth)

	if err := e.Start(":8080"); err != nil {
		glog.Fatal(err.Error())
	}
}

func checkHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, "healthy")
}

// handler to cleanup Gorilla Mux session context
func clearHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		context.Clear(c.Request())
		return next(c)
	}
}

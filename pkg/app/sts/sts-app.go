package sts

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/sts/handlers"
	"github.com/labstack/echo"

	mtMiddleware "github.com/kpn/pion/pkg/pion/shared/multi-tenancy/middleware"
)

type App struct {
}

func NewApp() App {
	return App{}
}

func (app App) Start() {
	glog.Info("Pion Security Token Service")

	e := echo.New()

	g := e.Group("/customers/:customerName")
	g.Use(mtMiddleware.ValidateCustomer)
	configureAccessKeyAPIs(g)

	// query secret key
	// TODO implement RPC for improve performance
	e.GET("/accesskey/:key", handlers.QueryAccessKey)

	e.GET("/health", checkHealth)

	if err := e.Start(":8080"); err != nil {
		glog.Fatal(err.Error())
	}
}

// configureAccessKeyAPIs provides APIs to generate a pair of accessKey/secretKey used for accessing Pion
func configureAccessKeyAPIs(g *echo.Group) {
	g.POST("/accesskey", handlers.CreateAccessKey)
	g.GET("/users/:username", handlers.ListAccessKeys)
	g.DELETE("/accesskey", handlers.RevokeAccessKey)
}

func checkHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, "healthy")
}

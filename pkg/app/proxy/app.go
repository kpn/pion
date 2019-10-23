package proxy

import (
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/proxy/handlers"
	"github.com/kpn/pion/pkg/pion/proxy/middleware"
	"github.com/kpn/pion/pkg/pion/proxy/watchers"
	"github.com/labstack/echo"
	"go.etcd.io/etcd/clientv3"
)

var DefaultEtcdAddress = os.Getenv("ETCD_ADDRESS")

type App struct {
}

// Verify Pion access keys and inject Minio access/secret keys
func NewApp() App {
	return App{}
}

func (app App) Start() {
	glog.Info("Pion Object Storage Service HandleRequest")
	app.MonitorPublicPaths()

	e := echo.New()

	e.GET("/*", handlers.HandleRequest, middleware.Authenticate, middleware.Authorize)
	e.HEAD("/*", handlers.HandleRequest, middleware.Authenticate, middleware.Authorize)
	e.POST("/*", handlers.HandleRequest, middleware.Authenticate, middleware.Authorize)
	e.PUT("/*", handlers.HandleRequest, middleware.Authenticate, middleware.Authorize, middleware.PopulateBucket)
	e.DELETE("/*", handlers.HandleRequest, middleware.Authenticate, middleware.Authorize, middleware.PopulateBucket)

	e.GET("/health", checkHealth)

	if err := e.Start(":8080"); err != nil {
		glog.Fatal(err.Error())
	}
}

func (app App) MonitorPublicPaths() {
	glog.V(2).Info("Monitoring public paths")
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{DefaultEtcdAddress},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		glog.Errorf("Failed to connect to etcd '%s': %v", DefaultEtcdAddress, err)
	}

	middleware.PublicPathsWatcher = watchers.NewPublicPathsWatcher(client)
	err = middleware.PublicPathsWatcher.Init()
	if err != nil {
		glog.Errorf("Failed to initialize public path monitoring: %v", err)
		return
	}
	middleware.PublicPathsWatcher.Watch()
}

func checkHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, "healthy")
}

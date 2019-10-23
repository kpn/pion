package authz

import (
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/authz/policy"
	"github.com/kpn/pion/pkg/pion/authz/watcher"
	"github.com/kpn/pion/pkg/sts/cache"
	"github.com/labstack/echo"
)

type App struct {
	engine         policy.Engine
	bucketsWatcher watcher.BucketsWatcher
}

var DefaultEtcdAddress = os.Getenv("ETCD_ADDRESS")

const defaultBucketPathPrefix = "/pion/oss/buckets"

func NewApp() App {
	return App{}
}

func (app *App) Start() {
	glog.Info("Pion Object Storage Service Authorization")

	app.initBucketsWatcher()

	engine, err := policy.NewEngine(app.bucketsWatcher)
	if err != nil {
		glog.Fatalf("Failed to load policies engine: %v", err)
	}
	app.engine = engine

	e := echo.New()

	e.GET("/data-authorize", app.DataPlaneAuthorize)
	e.GET("/mt-rbac-authorize", app.ControlPlaneAuthorize)
	e.GET("/health", checkHealth)

	if err := e.Start(":8080"); err != nil {
		glog.Fatal(err.Error())
	}
}

func (app *App) initBucketsWatcher() {
	ec, err := cache.NewEtcdClient(DefaultEtcdAddress)
	if err != nil {
		glog.Fatalf("Failed to connect to etcd service '%s': %v", DefaultEtcdAddress, err)
	}
	app.bucketsWatcher, err = watcher.NewBucketsWatcher(defaultBucketPathPrefix, ec)
	if err != nil {
		glog.Fatalf("Failed to watch buckets at '%s': %v", defaultBucketPathPrefix, err)
	}
	app.bucketsWatcher.Watch()
}

func checkHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, "healthy")
}

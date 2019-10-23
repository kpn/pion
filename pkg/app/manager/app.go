package manager

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/manager/handlers"
	"github.com/kpn/pion/pkg/pion/manager/middleware"
	mtMiddleware "github.com/kpn/pion/pkg/pion/shared/multi-tenancy/middleware"
	"github.com/kpn/pion/pkg/pion/shared/rbac"
	"github.com/labstack/echo"
)

type App struct {
}

func NewApp() App {
	return App{}
}

// TODO Implement leader election to support running in HA mode
func (app App) Start(etcdAddress string) {
	glog.Info("Pion Manager")

	e := echo.New()

	group := e.Group("/customers/:customerName")
	group.Use(mtMiddleware.ValidateCustomer, mtMiddleware.ValidateUserInCustomer)
	configureCustomerAPIs(group)

	internalGroup := e.Group("/_internal")
	configureInternalAPIs(internalGroup)

	e.GET("/health", checkHealth)

	if err := e.Start(":8080"); err != nil {
		glog.Fatal(err.Error())
	}
}

func configureInternalAPIs(g *echo.Group) {
	// TODO authorize who can call customers APIs
	// internal customer management APIs
	g.GET("/customers/:name", handlers.GetCustomer)
	g.GET("/customers", handlers.ListCustomers)
	g.POST("/customers", handlers.AddCustomer)
	g.PUT("/customers/:name", handlers.UpdateCustomer)
	g.DELETE("/customers", handlers.DeleteCustomer)
}

func configureCustomerAPIs(g *echo.Group) {
	g.GET("/public-objects", handlers.ListPublicObjects)
	g.POST("/public-objects",
		middleware.Authorize(rbac.BucketResource, rbac.Publish, handlers.AddPublicObject))
	g.DELETE("/public-objects",
		middleware.Authorize(rbac.BucketResource, rbac.Unpublish, handlers.DeletePublicObject))

	g.GET("/buckets",
		middleware.Authorize(rbac.BucketResource, rbac.List, handlers.ListBuckets))
	g.GET("/buckets/:name",
		middleware.Authorize(rbac.BucketResource, rbac.Get, handlers.GetBucket))
	g.POST("/buckets",
		middleware.Authorize(rbac.BucketResource, rbac.Create, handlers.AddBucket))
	g.DELETE("/buckets",
		middleware.Authorize(rbac.BucketResource, rbac.Delete, handlers.DeleteBucket))

	// ACL APIs
	g.GET("/buckets/:name/acl",
		middleware.Authorize(rbac.BucketAclResource, rbac.List, handlers.GetBucketACL))
	g.PUT("/buckets/:name/acl",
		middleware.Authorize(rbac.BucketAclResource, rbac.Update, handlers.PutBucketACL))
	g.DELETE("/buckets/:name/acl",
		middleware.Authorize(rbac.BucketAclResource, rbac.Delete, handlers.DeleteBucketACL))

	g.GET("/role-bindings/:name",
		middleware.Authorize(rbac.RoleBindingResource, rbac.Get, handlers.GetRoleBinding))
	g.GET("/role-bindings",
		middleware.Authorize(rbac.RoleBindingResource, rbac.List, handlers.ListRoleBindings))
	g.POST("/role-bindings",
		middleware.Authorize(rbac.RoleBindingResource, rbac.Create, handlers.AddRoleBinding))
	g.DELETE("/role-bindings",
		middleware.Authorize(rbac.RoleBindingResource, rbac.Delete, handlers.DeleteRoleBinding))
	// TODO Update role-binding API

	g.GET("/roles",
		middleware.Authorize(rbac.RoleResource, rbac.List, handlers.ListRoles))
}

func checkHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, "healthy")
}

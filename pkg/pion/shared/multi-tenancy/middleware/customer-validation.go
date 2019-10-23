package middleware

import (
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/etcd-store"
	"github.com/kpn/pion/pkg/sts/cache"
	"github.com/labstack/echo"
)

var defaultEtcdAddress = os.Getenv("ETCD_ADDRESS")

// ValidateCustomer verifies if the request referring to the valid customer in the system
func ValidateCustomer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		customerName := c.Param("customerName")
		if customerName == "" {
			return c.JSON(http.StatusBadRequest, "Customer not found")
		}

		etcdClient, err := cache.NewEtcdClient(defaultEtcdAddress)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to etcd")
		}
		defer cache.SilentClose(etcdClient)

		store := etcd_store.NewCustomerStore(shared.DefaultKeyPrefix, etcdClient)
		customer, err := store.Get(customerName)
		if err != nil {
			glog.Warningf("Requesting to non-existent customer: %v", err)
			return echo.NewHTTPError(http.StatusForbidden, "Customer not found")
		}
		c.Set(shared.CustomerKey, customer)
		return next(c)
	}
}

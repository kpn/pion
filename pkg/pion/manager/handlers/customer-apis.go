package handlers

import (
	"net/http"

	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/etcd-store"
	"github.com/kpn/pion/pkg/sts/cache"
	"github.com/labstack/echo"
)

// ListCustomers is the boilerplate function calling CustomerHandler::List
func ListCustomers(c echo.Context) error {
	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(etcdClient)

	cs := etcd_store.NewCustomerStore(shared.DefaultKeyPrefix, etcdClient)
	return NewCustomerHandler(cs).List(c)
}

// AddCustomer is the boilerplate function calling CustomerHandler::Add
func AddCustomer(c echo.Context) error {
	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(etcdClient)

	cs := etcd_store.NewCustomerStore(shared.DefaultKeyPrefix, etcdClient)
	return NewCustomerHandler(cs).Add(c)
}

// UpdateCustomer is the boilerplate function calling CustomerHandler::Update
func UpdateCustomer(c echo.Context) error {
	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(etcdClient)

	cs := etcd_store.NewCustomerStore(shared.DefaultKeyPrefix, etcdClient)
	return NewCustomerHandler(cs).Update(c)
}

// DeleteCustomer is the boilerplate function calling CustomerHandler::Delete
func DeleteCustomer(c echo.Context) error {
	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(etcdClient)

	cs := etcd_store.NewCustomerStore(shared.DefaultKeyPrefix, etcdClient)

	return NewCustomerHandler(cs).Delete(c)
}

// GetCustomer is the boilerplate function calling CustomerHandler::Get
func GetCustomer(c echo.Context) error {
	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(etcdClient)

	cs := etcd_store.NewCustomerStore(shared.DefaultKeyPrefix, etcdClient)

	return NewCustomerHandler(cs).Get(c)
}

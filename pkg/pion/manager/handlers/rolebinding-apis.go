package handlers

import (
	"errors"
	"net/http"
	"path"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/kpn/pion/pkg/pion/shared/rbac/etcd-store"
	"github.com/kpn/pion/pkg/sts/cache"
	"github.com/labstack/echo"
	"go.etcd.io/etcd/clientv3"
)

// ListRoleBindings is the boilerplate function calling RoleBindingHandler::List
func ListRoleBindings(c echo.Context) error {
	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(etcdClient)

	rbs, err := createRoleBindingStore(c, etcdClient)
	if err != nil {
		glog.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return NewRoleBindingHandler(rbs).List(c)
}

// AddRoleBinding is the boilerplate function calling RoleBindingHandler::Add
func AddRoleBinding(c echo.Context) error {
	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(etcdClient)

	rbs, err := createRoleBindingStore(c, etcdClient)
	if err != nil {
		glog.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return NewRoleBindingHandler(rbs).Add(c)
}

// DeleteRoleBinding is the boilerplate function calling RoleBindingHandler::Delete
func DeleteRoleBinding(c echo.Context) error {
	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(etcdClient)

	rbs, err := createRoleBindingStore(c, etcdClient)
	if err != nil {
		glog.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return NewRoleBindingHandler(rbs).Delete(c)
}

// GetRoleBinding is the boilerplate function calling RoleBindingHandler::Get
func GetRoleBinding(c echo.Context) error {
	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(etcdClient)

	rbs, err := createRoleBindingStore(c, etcdClient)
	if err != nil {
		glog.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return NewRoleBindingHandler(rbs).Get(c)
}

func createRoleBindingStore(c echo.Context, etcdClient *clientv3.Client) (etcd_store.RoleBindingStore, error) {
	customer, ok := c.Get(shared.CustomerKey).(*model.Customer)
	if !ok || customer == nil {
		return nil, errors.New("customer not found in the request context")
	}

	return etcd_store.NewRoleBindingStore(path.Join(shared.DefaultKeyPrefix, "rolebindings", customer.Name), etcdClient), nil
}

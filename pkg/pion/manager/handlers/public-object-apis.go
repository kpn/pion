package handlers

import (
	"errors"
	"net/http"
	"path"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/etcd-store"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/labstack/echo"
)

// ListPublicObjects is the boilerplate function calling PublicObjectHandler::List
func ListPublicObjects(c echo.Context) error {
	publicObjectStore, err := createPublicObjectStore(c)
	if err != nil {
		glog.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	h := NewPublicObjectHandler(publicObjectStore)
	return h.List(c)
}

// AddPublicObject is the boilerplate function calling PublicObjectHandler::Add
func AddPublicObject(c echo.Context) error {
	publicObjectStore, err := createPublicObjectStore(c)
	if err != nil {
		glog.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	h := NewPublicObjectHandler(publicObjectStore)
	return h.Add(c)
}

// DeletePublicObject is the boilerplate function calling PublicObjectHandler::Delete
func DeletePublicObject(c echo.Context) error {
	publicObjectStore, err := createPublicObjectStore(c)
	if err != nil {
		glog.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	h := NewPublicObjectHandler(publicObjectStore)
	return h.Delete(c)
}

func createPublicObjectStore(c echo.Context) (etcd_store.FileObjectStore, error) {
	customer, ok := c.Get(shared.CustomerKey).(*model.Customer)
	if !ok || customer == nil {
		return nil, errors.New("customer not found in the request context")
	}

	return etcd_store.NewFileObjectStore(path.Join(shared.DefaultKeyPrefix, "public-objects", customer.Name), shared.DefaultEtcdAddress)
}

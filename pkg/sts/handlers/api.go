package handlers

import (
	"net/http"
	"path"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/kpn/pion/pkg/sts/cache"
	store "github.com/kpn/pion/pkg/sts/pion-store"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	etcd "go.etcd.io/etcd/clientv3"
)

// CreateAccessKey is the boilerplate code to call AccessKeyHandler::Create
func CreateAccessKey(c echo.Context) error {
	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(etcdClient)

	ks, err := createEtcdKeyStoreByCustomer(c, etcdClient)
	if err != nil {
		glog.Error(err.Error())
		return Response(c, http.StatusInternalServerError, "Failed to create Etcd key store")
	}

	return NewAccessKeyHandler(ks).Create(c)
}

// ListAccessKeys is the boilerplate code to call AccessKeyHandler::List
func ListAccessKeys(c echo.Context) error {
	ec, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(ec)

	ks, err := createEtcdKeyStoreByCustomer(c, ec)
	if err != nil {
		glog.Error(err.Error())
		return Response(c, http.StatusInternalServerError, "Failed to create Etcd key store")
	}

	return NewAccessKeyHandler(ks).List(c)
}

// RevokeAccessKey is the boilerplate code to call AccessKeyHandler::Revoke
func RevokeAccessKey(c echo.Context) error {
	ec, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(ec)

	ks, err := createEtcdKeyStoreByCustomer(c, ec)
	if err != nil {
		glog.Error(err.Error())
		return Response(c, http.StatusInternalServerError, "Failed to create Etcd key store")
	}

	return NewAccessKeyHandler(ks).Revoke(c)
}

// QueryAccessKey handles API to query secret data from an access key
func QueryAccessKey(c echo.Context) error {
	accessKey := c.Param("key")
	if accessKey == "" {
		return c.JSON(http.StatusBadRequest, "Missing accessKey")
	}

	ec, err := cache.NewEtcdClient(cache.DefaultEtcdAddress)
	if err != nil {
		glog.Errorf("Failed to connect to Etcd: %v", err)
		return Response(c, http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(ec)

	kq, err := store.NewKeyQuerier(store.PayloadKeyPrefix, ec)
	if err != nil {
		return Response(c, http.StatusInternalServerError, "Failed to create key querier")
	}

	secretData, err := kq.Query(accessKey)
	if err != nil {
		glog.Errorf("Failed to get secret key: %v", err)
		return Response(c, http.StatusInternalServerError, "Failed to get secret key")
	}
	return c.JSON(http.StatusOK, secretData)
}

func createEtcdKeyStoreByCustomer(c echo.Context, client *etcd.Client) (ks store.KeyStore, err error) {
	customer, ok := c.Get(shared.CustomerKey).(*model.Customer)
	if !ok || customer == nil {
		return ks, errors.New("no customer found in the request context")
	}
	// Indices are clustered by customer for isolation, payloads are flat for fast querying from proxy
	customerIndexKeyPrefix := path.Join(store.IndexKeyPrefix, customer.Name) // grouping users' tokens by customer name
	return store.NewEtcdKeyStore(customerIndexKeyPrefix, store.PayloadKeyPrefix, client)
}

package handlers

import (
	"errors"
	"net/http"
	"path"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/etcd-store"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/validator"
	"github.com/kpn/pion/pkg/sts/cache"
	"github.com/labstack/echo"
	"go.etcd.io/etcd/clientv3"
)

// ListBuckets is the boilerplate function calling BucketHandler::List
func ListBuckets(c echo.Context) error {
	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(etcdClient)

	bktStore, err := createBucketStore(c, etcdClient)
	if err != nil {
		glog.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return NewBucketHandler(bktStore, createBucketValidator(etcdClient)).List(c)
}

// GetBucket is the boilerplate function calling BucketHandler::Get
func GetBucket(c echo.Context) error {
	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(etcdClient)

	bktStore, err := createBucketStore(c, etcdClient)
	if err != nil {
		glog.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return NewBucketHandler(bktStore, createBucketValidator(etcdClient)).Get(c)
}

// AddBucket is the boilerplate function calling BucketHandler::Add
func AddBucket(c echo.Context) error {
	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(etcdClient)

	bktStore, err := createBucketStore(c, etcdClient)
	if err != nil {
		glog.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return NewBucketHandler(bktStore, createBucketValidator(etcdClient)).Add(c)
}

// DeleteBucket is the boilerplate function calling BucketHandler::Delete
func DeleteBucket(c echo.Context) error {
	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to Etcd")
	}
	defer cache.SilentClose(etcdClient)

	bktStore, err := createBucketStore(c, etcdClient)
	if err != nil {
		glog.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return NewBucketHandler(bktStore, createBucketValidator(etcdClient)).Delete(c)
}

func createBucketStore(c echo.Context, etcdClient *clientv3.Client) (etcd_store.BucketStore, error) {
	customer, ok := c.Get(shared.CustomerKey).(*model.Customer)
	if !ok || customer == nil {
		return nil, errors.New("customer not found in the request context")
	}

	return etcd_store.NewBucketStore(path.Join(shared.DefaultKeyPrefix, "buckets", customer.Name), etcdClient), nil
}

func createBucketValidator(etcdClient *clientv3.Client) validator.BucketValidator {
	return validator.NewUniqueBucketValidator(path.Join(shared.DefaultKeyPrefix, "buckets"), etcdClient)
}

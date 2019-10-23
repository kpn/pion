package handlers

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/sts/cache"
	"github.com/labstack/echo"
)

// GetBucketACL is the boilerplate function calling BucketAclHandler::Get
func GetBucketACL(c echo.Context) error {
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
	return NewBucketACLHandler(bktStore).Get(c)
}

// PutBucketACL is the boilerplate function calling BucketAclHandler::Update
func PutBucketACL(c echo.Context) error {
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
	return NewBucketACLHandler(bktStore).Update(c)
}

// DeleteBucketACL is the boilerplate function calling BucketAclHandler::Delete
func DeleteBucketACL(c echo.Context) error {
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
	return NewBucketACLHandler(bktStore).Delete(c)
}

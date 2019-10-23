package handlers

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/manager/tenant"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/etcd-store"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/kpn/pion/pkg/pion/shared/s3"
	"github.com/kpn/pion/pkg/sts/cache"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type defaultPublicObjectHandler struct {
	fileObjectStore etcd_store.FileObjectStore
}

// NewPublicObjectHandler creates handler managing public-objects
func NewPublicObjectHandler(fos etcd_store.FileObjectStore) Handler {
	return &defaultPublicObjectHandler{
		fileObjectStore: fos,
	}
}

func (h defaultPublicObjectHandler) Get(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, "not implemented")
}

func (h defaultPublicObjectHandler) List(c echo.Context) error {
	paths, err := h.fileObjectStore.GetAllPaths()
	if err != nil {
		glog.Errorf("Failed to get public objects: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get public objects")
	}
	if len(paths) == 0 {
		glog.V(2).Info("No public path found")
	}
	glog.V(3).Infof("Public paths: %v", paths)
	return c.JSON(http.StatusOK, paths)
}

func (h defaultPublicObjectHandler) Add(c echo.Context) error {
	type Payload struct {
		Path string `json:"path"`
	}
	var p Payload
	err := c.Bind(&p)
	if err != nil {
		glog.Warningf("Payload is invalid or not found: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Payload is invalid or not found")
	}

	err = h.validate(c, p.Path)
	if err != nil {
		glog.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "invalid path")
	}
	existed, err := h.hasPathExisted(p.Path)
	if err != nil {
		glog.Errorf("Failed to check path existence: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add public object")
	}
	if existed {
		glog.Warningf("Public path '%s' already existed", p.Path)
		return echo.NewHTTPError(http.StatusBadRequest, "path already existed")
	}

	obj, err := h.fileObjectStore.Add(p.Path)
	if err != nil {
		glog.Errorf("Failed to add public object: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add public object")
	}
	glog.V(3).Infof("Added public paths: %v", obj)
	return c.JSON(http.StatusOK, obj)
}

func (h defaultPublicObjectHandler) Delete(c echo.Context) error {
	type Payload struct {
		Path string `json:"path"`
	}
	var p Payload
	err := c.Bind(&p)
	if err != nil {
		glog.Warningf("Payload is invalid or not found: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Payload is invalid or not found")
	}
	err = h.validate(c, p.Path)
	if err != nil {
		glog.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "invalid path")
	}

	existed, err := h.hasPathExisted(p.Path)
	if err != nil {
		glog.Errorf("Failed to check path existence: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add public object")
	}
	if !existed {
		return echo.NewHTTPError(http.StatusBadRequest, "path not existed")
	}

	err = h.fileObjectStore.Delete(p.Path)
	if err != nil {
		glog.Errorf("Failed to delete public object '%s': %v", p.Path, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete public object")
	}
	glog.V(3).Infof("Deleted public paths: %v", p.Path)
	return c.JSON(http.StatusOK, p)
}

func (h defaultPublicObjectHandler) Update(c echo.Context) error {
	return echo.NewHTTPError(http.StatusInternalServerError, "not implemented defaultPublicObjectHandler:Update")
}

func (h defaultPublicObjectHandler) validate(c echo.Context, publicPath string) error {
	if publicPath == "" {
		return errors.New("empty public path")
	}

	// check if the bucket is owned by the customer
	bucketName, _, err := s3.GetResources(publicPath)
	if err != nil {
		return err
	}

	return h.hasBucketExisted(c, bucketName)

}

func (h defaultPublicObjectHandler) hasBucketExisted(c echo.Context, bucketName string) error {
	customer, ok := c.Get(shared.CustomerKey).(*model.Customer)
	if !ok || customer == nil {
		return errors.New("customer not found in the context")
	}
	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		glog.Errorf("Failed to connect to etcd: %v", err)
		return err
	}
	defer cache.SilentClose(etcdClient)
	bucketStore := tenant.NewBucketStore(etcdClient, customer.Name)
	bkt, err := bucketStore.Get(bucketName)
	if err != nil {
		return err
	}
	if bkt == nil {
		return errors.New("Bucket not found")
	}
	return nil
}

func (h defaultPublicObjectHandler) hasPathExisted(path string) (bool, error) {
	allPaths, err := h.fileObjectStore.GetAllPaths()
	if err != nil {
		return false, err
	}
	for _, p := range allPaths {
		if p == path {
			return true, nil
		}
	}
	return false, nil
}

package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/authz"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/etcd-store"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/labstack/echo"
)

type bucketACLHandler struct {
	store etcd_store.BucketStore
}

// NewBucketACLHandler creates handler managing ACLs of buckets
func NewBucketACLHandler(s etcd_store.BucketStore) Handler {
	return &bucketACLHandler{
		store: s,
	}
}

func (h bucketACLHandler) List(c echo.Context) error {
	return echo.NewHTTPError(http.StatusInternalServerError, "not implemented bucketACLHandler:List")
}

func (h bucketACLHandler) Add(c echo.Context) error {
	return echo.NewHTTPError(http.StatusInternalServerError, "not implemented bucketACLHandler:Add")
}

func (h bucketACLHandler) Get(c echo.Context) error {
	bucket, err := h.getBucket(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, bucket.ACLs)
}

func (h bucketACLHandler) Update(c echo.Context) error {
	bktName := c.Param("name")
	if bktName == "" {
		return c.JSON(http.StatusBadRequest, "Missing bucket name param")
	}

	var acls authz.ACLList
	err := c.Bind(&acls)
	if err != nil {
		glog.Warningf("ACL payload is invalid or not found: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Payload is invalid or not found")
	}

	err = h.store.UpdateACLs(bktName, acls)
	if err != nil {
		glog.Errorf("Saving ACLs for bucket '%s' failed: %v", bktName, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Saving ACLs failed")
	}
	return c.JSON(http.StatusOK, acls)
}

func (h bucketACLHandler) Delete(c echo.Context) error {
	bktName := c.Param("name")
	if bktName == "" {
		return c.JSON(http.StatusBadRequest, "Missing bucket name param")
	}

	err := h.store.UpdateACLs(bktName, nil)
	if err != nil {
		glog.Errorf("Failed to delete ACLs of bucket '%s': %v", bktName, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete bucket ACLs")
	}
	return c.NoContent(http.StatusOK)
}

func (h bucketACLHandler) getBucket(c echo.Context) (*model.Bucket, error) {
	name := c.Param("name")
	if name == "" {
		return nil, c.JSON(http.StatusBadRequest, "Missing bucket name param")
	}
	glog.V(2).Infof("Getting bucket '%s'", name)

	bucket, err := h.store.Get(name)
	if err != nil {
		glog.Errorf("Failed to get bucket from DB: %v", err)
		return nil, c.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot get bucket '%s'", name))
	}
	return bucket, nil
}

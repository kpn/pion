package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/etcd-store"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/validator"
	"github.com/labstack/echo"
)

type bucketHandler struct {
	bucketValidator validator.BucketValidator
	store           etcd_store.BucketStore
}

// NewBucketHandler creates handler managing buckets
func NewBucketHandler(s etcd_store.BucketStore, v validator.BucketValidator) Handler {
	return &bucketHandler{
		store:           s,
		bucketValidator: v,
	}
}

func (h bucketHandler) Get(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, "Missing bucket name param")
	}
	glog.V(2).Infof("Getting bucket '%s'", name)

	rb, err := h.store.Get(name)
	if err != nil {
		glog.Errorf("Failed to get bucket from DB: %v", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot get bucket '%s'", name))
	}

	return c.JSON(http.StatusOK, rb)
}

func (h bucketHandler) Add(c echo.Context) error {
	var b model.Bucket
	err := c.Bind(&b)
	if err != nil {
		glog.Warningf("Payload is invalid or not found: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Payload is invalid or not found")
	}

	err = h.bucketValidator.Validate(b.Name)
	if err != nil {
		glog.Errorf("Validating bucket '%s' failed: %v", b.Name, err)
		switch err {
		case validator.ErrBucketExisted:
			return c.String(http.StatusConflict, err.Error())
		default:
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	// prevented overriding generated fields
	b.ACLs = nil
	b.Creator = ""
	b.CreatedAt = time.Now()

	h.setCreator(c, &b)
	h.setCustomer(c, &b)

	nb, err := h.store.Add(b)
	if err != nil {
		glog.Errorf("Failed to add new bucket: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add new bucket")
	}
	return c.JSON(http.StatusOK, nb)
}

func (h bucketHandler) Delete(c echo.Context) error {
	type Payload struct {
		Name string `json:"name"`
	}

	var p Payload
	err := c.Bind(&p)
	if err != nil {
		glog.Warningf("Payload is invalid or not found: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Payload is invalid or not found")
	}
	if p.Name == "" {
		glog.Warningf("Empty name in the payload")
		return echo.NewHTTPError(http.StatusBadRequest, "Empty name in the payload")
	}

	err = h.store.Delete(p.Name)
	if err != nil {
		glog.Errorf("Failed to delete bucket '%s': %v", p.Name, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete bucket")
	}
	return c.NoContent(http.StatusOK)
}

func (h bucketHandler) List(c echo.Context) error {
	buckets, err := h.store.List()
	if err != nil {
		glog.Errorf("Failed to list buckets: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list buckets")
	}

	glog.V(2).Infof("Listed buckets: %v", buckets)
	return c.JSON(http.StatusOK, buckets)
}

func (h bucketHandler) Update(c echo.Context) error {
	return echo.NewHTTPError(http.StatusInternalServerError, "not implemented bucketHandler:Update")
}

func (h bucketHandler) setCustomer(c echo.Context, b *model.Bucket) {
	customer, ok := c.Get(shared.CustomerKey).(*model.Customer)
	if !ok {
		glog.Errorf("Failed to get customer from request context")
	}
	b.OwnedBy = customer.Name
}

func (h bucketHandler) setCreator(c echo.Context, b *model.Bucket) {
	userId := c.Request().Header.Get(shared.UserIdKey)
	if userId == "" {
		glog.Errorf("Header value of '%s' not found", shared.UserIdKey)
	} else {
		b.Creator = userId
	}
}

package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/etcd-store"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/labstack/echo"
)

type customerHandler struct {
	store etcd_store.CustomerStore
}

// NewCustomerHandler creates handler managing customers
func NewCustomerHandler(s etcd_store.CustomerStore) Handler {
	return &customerHandler{
		store: s,
	}
}

func (h customerHandler) Get(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, "Missing customer name param")
	}
	glog.V(2).Infof("Getting customer '%s'", name)

	customer, err := h.store.Get(name)
	if err != nil {
		glog.Warningf("Failed to get customer from DB: %v", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot get customer '%s'", name))
	}

	return c.JSON(http.StatusOK, customer)
}

func (h customerHandler) Add(c echo.Context) error {
	var customer model.Customer

	err := c.Bind(&customer)
	if err != nil {
		glog.Warningf("Payload is invalid or not found: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Payload is invalid or not found")
	}

	newCustomer, err := h.store.Add(customer)
	if err != nil {
		glog.Errorf("Failed to add new customer: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add new customer")
	}
	glog.V(1).Infof("Added customer '%s'", customer.Name)
	return c.JSON(http.StatusOK, newCustomer)
}

func (h customerHandler) Delete(c echo.Context) error {
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
		glog.Errorf("Failed to delete customer '%s': %v", p.Name, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete customer")
	}
	glog.V(1).Infof("Deleted customer '%s'", p.Name)
	return c.NoContent(http.StatusOK)
}

func (h customerHandler) List(c echo.Context) error {
	customers, err := h.store.List()
	if err != nil {
		glog.Warningf("Failed to list customers: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list customers")
	}

	glog.V(2).Infof("Listed customers: %v", customers)
	return c.JSON(http.StatusOK, customers)
}

func (h customerHandler) Update(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, "Missing customer name param")
	}

	var payload model.Customer
	err := c.Bind(&payload)
	if err != nil {
		glog.Warningf("Payload is invalid or not found: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Payload is invalid or not found")
	}

	customer, err := h.store.Get(name)
	if err != nil || customer == nil {
		glog.Warningf("Customer '%s' not found", name)
		if err != nil {
			glog.Warningf("Error: %v", err)
		}
		return echo.NewHTTPError(http.StatusBadRequest, "customer not found")
	}
	// only allow to update groups and userIDs attributes
	customer.Groups = payload.Groups
	customer.UserIDs = payload.UserIDs

	updatedCustomer, err := h.store.Update(*customer)
	if err != nil {
		glog.Errorf("Failed to update customer: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update customer")
	}
	glog.V(1).Infof("Updated customer '%s'", customer.Name)
	return c.JSON(http.StatusOK, updatedCustomer)
}

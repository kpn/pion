package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared/rbac"
	"github.com/kpn/pion/pkg/pion/shared/rbac/etcd-store"
	"github.com/labstack/echo"
)

type roleBindingHandler struct {
	store etcd_store.RoleBindingStore
}

// NewRoleBindingHandler creates handler managing role-binding objects
func NewRoleBindingHandler(s etcd_store.RoleBindingStore) Handler {
	return &roleBindingHandler{
		store: s,
	}
}

func (h roleBindingHandler) Get(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, "Missing role-binding name param")
	}
	glog.V(2).Infof("Getting role-binding '%s'", name)

	rb, err := h.store.Get(name)
	if err != nil {
		glog.Errorf("Failed to get role-binding from DB: %v", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot get role-binding '%s'", name))
	}

	return c.JSON(http.StatusOK, rb)
}

func (h roleBindingHandler) Add(c echo.Context) error {
	var rb rbac.RoleBinding

	err := c.Bind(&rb)
	if err != nil {
		glog.Warningf("Payload is invalid or not found: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Payload is invalid or not found")
	}

	newRb, err := h.store.Add(rb)
	if err != nil {
		glog.Errorf("Failed to add new role-binding: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add new role-binding")
	}
	return c.JSON(http.StatusOK, newRb)
}

func (h roleBindingHandler) Delete(c echo.Context) error {
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
		glog.Errorf("Failed to delete role-binding '%s': %v", p.Name, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete role-binding")
	}
	return c.NoContent(http.StatusOK)
}

func (h roleBindingHandler) List(c echo.Context) error {
	rbs, err := h.store.List()
	if err != nil {
		glog.Errorf("Failed to list role-bindings: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list role-bindings")
	}

	glog.V(2).Infof("Listed role-bindings: %v", rbs)
	return c.JSON(http.StatusOK, rbs)
}

func (h roleBindingHandler) Update(c echo.Context) error {
	return echo.NewHTTPError(http.StatusInternalServerError, "not implemented")
}

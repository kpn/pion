package handlers

import (
	"net/http"
	"path"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/authz/request"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/rbac"
	"github.com/kpn/pion/pkg/pion/shared/rbac/etcd-store"
	"github.com/kpn/pion/pkg/sts/cache"
)

// RBACHandler defines APIs to handle RBAC authorization
type RBACHandler interface {
	Authorize(rbacRequest *rbac.Request) rbac.Decision

	AuthorizeHTTP(r *http.Request) rbac.Decision
}

type rbacHandler struct {
	rbacEngine rbac.Engine
}

// NewRBACHandler loads role-bindings of the specific customer
func NewRBACHandler(etcdAddress string, customerName string) (RBACHandler, error) {
	etcdClient, err := cache.NewEtcdClient(etcdAddress)
	if err != nil {
		return nil, err
	}
	defer cache.SilentClose(etcdClient)

	// system role-bindings
	systemRBStore := etcd_store.NewRoleBindingStore(path.Join(shared.DefaultKeyPrefix, "rolebindings", shared.SystemKey), etcdClient)
	systemRBs, err := systemRBStore.List()
	if err != nil {
		glog.Errorf("Failed to query role-binding of the system: %v", err)
		return nil, err
	}
	glog.V(3).Infof("System role-bindings: '%+v'", systemRBs)

	// customer role-bindings
	customerRBStore := etcd_store.NewRoleBindingStore(path.Join(shared.DefaultKeyPrefix, "rolebindings", customerName), etcdClient)
	customerRBs, err := customerRBStore.List()
	if err != nil {
		glog.Errorf("Failed to query role-binding of customer '%s': %v", customerName, err)
		return nil, err
	}
	glog.V(3).Infof("Customer role-bindings: '%+v'", customerRBs)

	allRBs := append(systemRBs, customerRBs...)
	engine := rbac.NewEngine(rbac.DefaultRoles, allRBs)

	return &rbacHandler{
		rbacEngine: engine,
	}, nil
}

func (handler rbacHandler) AuthorizeHTTP(r *http.Request) rbac.Decision {
	userId, userGroups := request.GetUserAttributes(r)
	if userId == "" || len(userGroups) == 0 {
		glog.Warningf("Invalid request")
		return rbac.DenyDecision
	}
	action, err := rbac.Str2Action(r.Header.Get(shared.ActionKey))
	if err != nil {
		glog.Error(err)
		return rbac.DenyDecision
	}

	resource, err := rbac.Str2Resource(r.Header.Get(shared.ResourceKey))
	if err != nil {
		glog.Error(err)
		return rbac.DenyDecision
	}

	azr := &rbac.Request{
		UserID: userId,
		Groups: userGroups,
		Target: resource,
		Action: action,
	}

	return handler.Authorize(azr)
}

func (handler rbacHandler) Authorize(rbacRequest *rbac.Request) rbac.Decision {
	glog.V(3).Infof("Evaluating request '%v", rbacRequest)
	// evaluate
	return handler.rbacEngine.Evaluate(*rbacRequest)
}

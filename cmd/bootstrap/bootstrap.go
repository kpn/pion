package main

import (
	"path"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/rbac"
	rbac_store "github.com/kpn/pion/pkg/pion/shared/rbac/etcd-store"
	"github.com/kpn/pion/pkg/sts/cache"
)

func initDefaultSuperAdmin() {
	var defaultRoleBinding = rbac.RoleBinding{
		Name:    "default-admin-rb",
		RoleRef: "admin",
		Subjects: []rbac.Subject{
			{Type: rbac.UserType, Value: adminUser},
		},
	}

	glog.Infof("Adding '%s' as the Admin of the system", adminUser)

	etcdClient, err := cache.NewEtcdClient(shared.DefaultEtcdAddress)
	if err != nil {
		glog.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer cache.SilentClose(etcdClient)

	glog.Info("Adding system level role-binding")
	roleBindingStore := rbac_store.NewRoleBindingStore(path.Join(shared.DefaultKeyPrefix, "rolebindings", shared.SystemKey), etcdClient)
	existedRB, err := roleBindingStore.Get(defaultRoleBinding.Name)
	if err == nil && existedRB != nil {
		glog.Infof("Role-binding '%s' existed, deleting", defaultRoleBinding.Name)
		err = roleBindingStore.Delete(defaultRoleBinding.Name)
		if err != nil {
			glog.Errorf("Failed to delete existing role-binding: %v", err)
		}
	}

	_, err = roleBindingStore.Add(defaultRoleBinding)
	if err != nil {
		glog.Errorf("Adding default system role-bindings failed: %v", err)
	}
	glog.Infof("Created role-binding: %v", defaultRoleBinding)
}

package etcd_store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared/rbac"
	etcd "go.etcd.io/etcd/clientv3"
)

// RoleBindingStore defines API to manage role-binding objects
type RoleBindingStore interface {
	Add(rb rbac.RoleBinding) (*rbac.RoleBinding, error)

	Delete(name string) error

	Get(name string) (*rbac.RoleBinding, error)

	List() ([]rbac.RoleBinding, error)
}

type roleBindingStore struct {
	keyPrefix string
	client    *etcd.Client
}

// NewRoleBindingStore creates store to manage role-binding in Etcd
func NewRoleBindingStore(keyPrefix string, client *etcd.Client) RoleBindingStore {
	return roleBindingStore{
		keyPrefix: keyPrefix,
		client:    client,
	}
}

func (store roleBindingStore) Add(rb rbac.RoleBinding) (*rbac.RoleBinding, error) {
	glog.V(2).Infof("adding role-binding %v", rb)
	indexKey := store.generateIndexKey(rb.Name)
	if store.exists(indexKey) {
		return nil, fmt.Errorf("role-binding '%s' already existed", rb.Name)
	}
	rb.CreatedAt = time.Now().UTC()

	bytes, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}

	_, err = store.client.Put(context.Background(), indexKey, string(bytes))
	if err != nil {
		glog.Errorf("Failed to save new role-binding '%s': %v", rb.Name, err)
		return nil, err
	}
	glog.V(1).Infof("Added new role-binding '%s'", rb.Name)
	return &rb, nil
}

func (store roleBindingStore) Delete(name string) error {
	indexKey := store.generateIndexKey(name)
	resp, err := store.client.Delete(context.Background(), indexKey)
	if err != nil {
		glog.Errorf("Failed to delete customer '%s': %v", name, err)
		return err
	}
	if resp.Deleted != 1 {
		err = fmt.Errorf("key '%s' not found", indexKey)
		glog.Error(err.Error())
		return err
	}
	glog.V(1).Infof("Deleted role-binding '%s'", name)
	return nil
}

func (store roleBindingStore) Get(name string) (*rbac.RoleBinding, error) {
	indexKey := store.generateIndexKey(name)
	resp, err := store.client.Get(context.Background(), indexKey)
	if err != nil {
		glog.Errorf("Failed to get etcd key '%s': %v", indexKey, err)
		return nil, err
	}

	if len(resp.Kvs) < 1 {
		glog.Errorf("Reading Etcd K-V error at '%s': should have 1 key-value", indexKey)
		return nil, errors.New("reading error: K-V does not exist")
	}
	var obj rbac.RoleBinding
	err = json.Unmarshal(resp.Kvs[0].Value, &obj)
	if err != nil {
		glog.Errorf("Failed to unmarshal role-binding object: %v", err)
		return nil, err
	}
	glog.V(2).Infof("Get bucket '%s: %v'", name, obj)
	return &obj, nil
}

func (store roleBindingStore) List() (bindings []rbac.RoleBinding, err error) {
	indexKeyPrefix := store.generateIndexKey("")
	resp, err := store.client.Get(context.Background(), indexKeyPrefix, etcd.WithPrefix())
	if err != nil {
		glog.Errorf("Failed to list Etcd keys with prefix '%s': %v", indexKeyPrefix, err)
		return nil, err
	}
	for _, kv := range resp.Kvs {
		var obj rbac.RoleBinding
		err = json.Unmarshal(kv.Value, &obj)
		if err != nil {
			glog.Errorf("Failed to unmarshal role-binding object: %v", err)
			continue
		}
		bindings = append(bindings, obj)
	}
	return
}

func (store roleBindingStore) generateIndexKey(name string) string {
	if name == "" {
		return store.keyPrefix + "/"
	}

	return path.Join(store.keyPrefix, name)
}

func (store roleBindingStore) exists(indexKey string) bool {
	resp, err := store.client.Get(context.Background(), indexKey, etcd.WithCountOnly())
	if err != nil {
		glog.Errorf("Failed to check key existence of '%s': %v", indexKey, err)
		return false
	}
	glog.V(2).Infof("Returned %d keys at '%s", resp.Count, indexKey)
	return resp.Count == 1
}

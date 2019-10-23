package etcd_store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	etcd "go.etcd.io/etcd/clientv3"
)

// CustomerStore defines API to manage customer store
type CustomerStore interface {
	Add(customer model.Customer) (*model.Customer, error)

	Update(customer model.Customer) (*model.Customer, error)

	Delete(name string) error

	Get(name string) (*model.Customer, error)

	List() ([]model.Customer, error)
}

type customerStore struct {
	keyPrefix string
	client    *etcd.Client
}

// NewCustomerStore creates new store object managing customers in Etcd
func NewCustomerStore(keyPrefix string, client *etcd.Client) CustomerStore {
	return customerStore{
		keyPrefix: keyPrefix,
		client:    client,
	}
}

func (store customerStore) Add(customer model.Customer) (*model.Customer, error) {
	indexKey := store.generateIndexKey(customer.Name)
	if store.exists(indexKey) {
		return nil, fmt.Errorf("customer '%s' already existed", customer.Name)
	}
	customer.CreatedAt = time.Now().UTC()
	customer.ModifiedAt = customer.CreatedAt

	return store.upsert(indexKey, customer)
}

func (store customerStore) Update(customer model.Customer) (*model.Customer, error) {
	indexKey := store.generateIndexKey(customer.Name)
	if !store.exists(indexKey) {
		return nil, fmt.Errorf("customer '%s' does not existed", customer.Name)
	}
	customer.ModifiedAt = time.Now().UTC()

	return store.upsert(indexKey, customer)
}

func (store customerStore) Delete(name string) error {
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
	return nil
}

func (store customerStore) Get(name string) (*model.Customer, error) {
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
	var obj model.Customer
	err = json.Unmarshal(resp.Kvs[0].Value, &obj)
	if err != nil {
		glog.Errorf("Failed to unmarshal customer object: %v", err)
		return nil, err
	}
	return &obj, nil
}

func (store customerStore) List() (customers []model.Customer, err error) {
	indexKeyPrefix := store.generateIndexKey("")
	resp, err := store.client.Get(context.Background(), indexKeyPrefix, etcd.WithPrefix())
	if err != nil {
		glog.Errorf("Failed to list Etcd keys with prefix '%s': %v", indexKeyPrefix, err)
		return nil, err
	}
	for _, kv := range resp.Kvs {
		var obj model.Customer
		err = json.Unmarshal(kv.Value, &obj)
		if err != nil {
			glog.Errorf("Failed to unmarshal customer object: %v", err)
			continue
		}
		glog.V(3).Infof("Got customer '%s'", obj.Name)
		customers = append(customers, obj)
	}
	return
}

func (store customerStore) generateIndexKey(name string) string {
	if name == "" {
		return path.Join(store.keyPrefix, "customers") + "/"
	}

	return path.Join(store.keyPrefix, "customers", name)
}

func (store customerStore) exists(indexKey string) bool {
	resp, err := store.client.Get(context.Background(), indexKey, etcd.WithCountOnly())
	if err != nil {
		glog.Errorf("Failed to check key existence of '%s': %v", indexKey, err)
		return false
	}
	glog.V(2).Infof("Returned %d keys at '%s", resp.Count, indexKey)
	return resp.Count == 1
}

func (store customerStore) upsert(indexKey string, customer model.Customer) (*model.Customer, error) {
	bytes, err := json.Marshal(customer)
	if err != nil {
		return nil, err
	}

	_, err = store.client.Put(context.Background(), indexKey, string(bytes))
	if err != nil {
		glog.Errorf("Failed to upsert customer '%s': %v", customer.Name, err)
		return nil, err
	}
	return &customer, nil
}
